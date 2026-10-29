package boardgamegeek

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL   = "https://boardgamegeek.com/xmlapi2/"
	defaultUserAgent = "go-boardgamegeek"
	defaultRetryWait = 2 * time.Second
	maxRetryWait     = time.Minute
)

// Client makes requests to the BGG xmlapi2. Build one with NewClient.
type Client struct {
	httpClient *http.Client
	token      string
	baseURL    *url.URL
	userAgent  string
	maxRetries int
	// retryWait is the first backoff delay, doubled on each retry
	retryWait time.Duration
}

// ClientOption configures a Client in NewClient
type ClientOption func(*Client)

// WithAuthToken sends token as a bearer token with every request
func WithAuthToken(token string) ClientOption {
	return func(c *Client) {
		c.token = token
	}
}

// WithHTTPClient sends requests with httpClient instead of a new http.Client,
// e.g. to set timeouts or a rate-limiting transport
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL replaces the xmlapi2 root that endpoint paths are joined to, e.g. to point at a test server
func WithBaseURL(baseURL *url.URL) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithUserAgent replaces the User-Agent header sent with every request
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithMaxRetries retries a request up to maxRetries times when BGG responds with 202 (request queued)
// or 429 (rate limited). It waits for Retry-After when BGG sends it, and otherwise backs off
// exponentially from 2s up to a minute. The default is 0: no retries.
func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// NewClient returns a Client for the BGG xmlapi2, configured by options
func NewClient(options ...ClientOption) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
		userAgent:  defaultUserAgent,
		retryWait:  defaultRetryWait,
	}
	for _, option := range options {
		option(c)
	}

	return c
}

// get requests path with opts, a request struct encoded by its `url` tags, and unmarshals the xml body into v.
// It retries 202s and 429s as configured by WithMaxRetries.
func (c *Client) get(ctx context.Context, path string, opts any, v any) error {
	params, encodeErr := query.Values(opts)
	if encodeErr != nil {
		return fmt.Errorf("error encoding %s query %w", path, encodeErr)
	}

	reqURL := c.baseURL.JoinPath(path)
	reqURL.RawQuery = params.Encode()

	for attempt := 0; ; attempt++ {
		err := c.do(ctx, reqURL, v)
		wait, retry := c.retryDelay(err, attempt)
		if !retry {
			return err
		}

		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}

// retryDelay reports whether a request that failed with err on attempt should be retried, and after how long
func (c *Client) retryDelay(err error, attempt int) (time.Duration, bool) {
	if attempt >= c.maxRetries {
		return 0, false
	}

	backoff := min(c.retryWait<<attempt, maxRetryWait)

	var rateLimitErr *RateLimitError
	if errors.As(err, &rateLimitErr) {
		if rateLimitErr.RetryAfter > 0 {
			return rateLimitErr.RetryAfter, true
		}
		return backoff, true
	}

	var acceptedErr *AcceptedError
	if errors.As(err, &acceptedErr) {
		return backoff, true
	}

	return 0, false
}

// do sends one request for reqURL and unmarshals the xml body into v.
// It returns *AcceptedError for 202, *RateLimitError for 429 and *ErrorResponse for any other
// non-200 status or for a 200 <error>/<errors> document or non-xml body. 404s and <error>
// documents match ErrNotFound; endpoints mark their own not-found quirks with markNotFound.
func (c *Client) do(ctx context.Context, reqURL *url.URL, v any) error {
	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if reqErr != nil {
		return fmt.Errorf("error creating %s request %w", reqURL.Path, reqErr)
	}

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, bggErr := c.httpClient.Do(req)
	if bggErr != nil {
		return fmt.Errorf("received error response from boardgamegeek %w", bggErr)
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusAccepted:
		return &AcceptedError{Response: resp}
	case http.StatusTooManyRequests:
		return &RateLimitError{Response: resp, RetryAfter: parseRetryAfter(resp.Header.Get("Retry-After"))}
	default:
		return &ErrorResponse{Response: resp, Message: http.StatusText(resp.StatusCode), notFound: resp.StatusCode == http.StatusNotFound}
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "xml") {
		return &ErrorResponse{Response: resp, Message: fmt.Sprintf("unexpected Content-Type %q", resp.Header.Get("Content-Type"))}
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return fmt.Errorf("error reading response body %w", readErr)
	}

	var bggErrResp apiError
	if xml.Unmarshal(body, &bggErrResp) == nil {
		return &ErrorResponse{Response: resp, Message: bggErrResp.Message, notFound: true}
	}

	var collectionErrs collectionErrors
	if xml.Unmarshal(body, &collectionErrs) == nil {
		return &ErrorResponse{Response: resp, Message: strings.Join(collectionErrs.Messages, "; ")}
	}

	if parseErr := xml.Unmarshal(body, v); parseErr != nil {
		return fmt.Errorf("error parsing xml body %w", parseErr)
	}

	return nil
}

// markNotFound makes err match ErrNotFound if it is an *ErrorResponse that notFound reports on
func markNotFound(err error, notFound func(*ErrorResponse) bool) error {
	var errResp *ErrorResponse
	if errors.As(err, &errResp) && notFound(errResp) {
		errResp.notFound = true
	}
	return err
}

// Search searches for items by name (/search)
func (c *Client) Search(ctx context.Context, opts *SearchOptions) (*SearchResponse, error) {
	var searchResponse SearchResponse
	if err := c.get(ctx, "/search", opts, &searchResponse); err != nil {
		return nil, err
	}
	return &searchResponse, nil
}

// GetThings returns one or more things such as board games, expansions or RPG items (/thing)
func (c *Client) GetThings(ctx context.Context, opts *ThingsOptions) (*ThingsResponse, error) {
	var getThingsResponse ThingsResponse
	if err := c.get(ctx, "/thing", opts, &getThingsResponse); err != nil {
		return nil, err
	}
	return &getThingsResponse, nil
}

// GetFamilyItems returns families, groups of related things such as a game series (/family)
func (c *Client) GetFamilyItems(ctx context.Context, opts *FamilyItemsOptions) (*FamilyItemsResponse, error) {
	var familyItemsResponse FamilyItemsResponse
	if err := c.get(ctx, "/family", opts, &familyItemsResponse); err != nil {
		return nil, err
	}
	return &familyItemsResponse, nil
}

// GetForumList returns the forums of a thing or family (/forumlist)
func (c *Client) GetForumList(ctx context.Context, opts *ForumListOptions) (*ForumListResponse, error) {
	var forumListResponse ForumListResponse
	if err := c.get(ctx, "/forumlist", opts, &forumListResponse); err != nil {
		return nil, err
	}
	return &forumListResponse, nil
}

// GetForum returns an error matching ErrNotFound for an unknown forum id
func (c *Client) GetForum(ctx context.Context, opts *ForumOptions) (*ForumResponse, error) {
	var forumResponse ForumResponse
	if err := c.get(ctx, "/forum", opts, &forumResponse); err != nil {
		// BGG responds to an unknown forum id with a 200 html page rather than xml
		err = markNotFound(err, func(e *ErrorResponse) bool {
			return e.Response.StatusCode == http.StatusOK && strings.Contains(e.Response.Header.Get("Content-Type"), "html")
		})
		return nil, fmt.Errorf("forum %d: %w", opts.ID, err)
	}
	return &forumResponse, nil
}

// GetThread returns an error matching ErrNotFound for an unknown thread id
func (c *Client) GetThread(ctx context.Context, opts *ThreadOptions) (*ThreadResponse, error) {
	var threadResponse ThreadResponse
	if err := c.get(ctx, "/thread", opts, &threadResponse); err != nil {
		return nil, fmt.Errorf("thread %d: %w", opts.ID, err)
	}
	return &threadResponse, nil
}

// GetHotItems returns the current hotness list, 50 items of one type (/hot)
func (c *Client) GetHotItems(ctx context.Context, opts *HotItemsOptions) (*HotItemsResponse, error) {
	var hotItemsResponse HotItemsResponse
	if err := c.get(ctx, "/hot", opts, &hotItemsResponse); err != nil {
		return nil, err
	}
	return &hotItemsResponse, nil
}

// GetUser returns an error matching ErrNotFound for an unknown username
func (c *Client) GetUser(ctx context.Context, opts *UserOptions) (*UserResponse, error) {
	var userResponse UserResponse
	if err := c.get(ctx, "/user", opts, &userResponse); err != nil {
		return nil, fmt.Errorf("user %s: %w", opts.Name, err)
	}
	return &userResponse, nil
}

// GetGuild returns an error matching ErrNotFound for an unknown guild id
func (c *Client) GetGuild(ctx context.Context, opts *GuildOptions) (*GuildResponse, error) {
	var guildResponse GuildResponse
	if err := c.get(ctx, "/guild", opts, &guildResponse); err != nil {
		return nil, fmt.Errorf("guild %d: %w", opts.ID, err)
	}
	return &guildResponse, nil
}

// GetCollection returns an *AcceptedError matching ErrCollectionNotReady while BGG is still building
// the collection; callers should retry the same request after a few seconds.
// An unknown username returns an error matching ErrNotFound.
func (c *Client) GetCollection(ctx context.Context, opts *CollectionOptions) (*CollectionResponse, error) {
	var collectionResponse CollectionResponse
	if err := c.get(ctx, "/collection", opts, &collectionResponse); err != nil {
		// BGG responds to an unknown username with a 200 <errors> document
		err = markNotFound(err, func(e *ErrorResponse) bool {
			return strings.Contains(e.Message, "Invalid username specified")
		})
		return nil, fmt.Errorf("collection for %s: %w", opts.Username, err)
	}
	return &collectionResponse, nil
}
