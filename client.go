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

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL   = "https://boardgamegeek.com/xmlapi2/"
	defaultUserAgent = "BoardGameGeekClient-go"
)

type Client struct {
	httpClient *http.Client
	token      string
	baseURL    *url.URL
	userAgent  string
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

// NewClient returns a Client for the BGG xmlapi2, configured by options
func NewClient(options ...ClientOption) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
		userAgent:  defaultUserAgent,
	}
	for _, option := range options {
		option(c)
	}

	return c
}

// get requests path with opts, a request struct encoded by its `url` tags, and unmarshals the xml body into v.
// It returns *AcceptedError for 202, *RateLimitError for 429 and *ErrorResponse for any other
// non-200 status or for a 200 <error>/<errors> document or non-xml body. 404s and <error>
// documents match ErrNotFound; endpoints mark their own not-found quirks with markNotFound.
func (bggClient *Client) get(ctx context.Context, path string, opts any, v any) error {
	params, encodeErr := query.Values(opts)
	if encodeErr != nil {
		return fmt.Errorf("error encoding %s query %w", path, encodeErr)
	}

	reqURL := bggClient.baseURL.JoinPath(path)
	reqURL.RawQuery = params.Encode()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if reqErr != nil {
		return fmt.Errorf("error creating %s request %w", path, reqErr)
	}

	if bggClient.userAgent != "" {
		req.Header.Set("User-Agent", bggClient.userAgent)
	}
	if bggClient.token != "" {
		req.Header.Set("Authorization", "Bearer "+bggClient.token)
	}

	resp, bggErr := bggClient.httpClient.Do(req)
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

func (bggClient *Client) Search(ctx context.Context, sr *SearchRequest) (*SearchResponse, error) {
	var searchResponse SearchResponse
	if err := bggClient.get(ctx, "/search", sr, &searchResponse); err != nil {
		return nil, err
	}
	return &searchResponse, nil
}

func (bggClient *Client) GetThings(ctx context.Context, gr *GetThingsRequest) (*GetThingsResponse, error) {
	var getThingsResponse GetThingsResponse
	if err := bggClient.get(ctx, "/thing", gr, &getThingsResponse); err != nil {
		return nil, err
	}
	return &getThingsResponse, nil
}

func (bggClient *Client) GetFamilyItems(ctx context.Context, fr *FamilyItemRequest) (*FamilyItemsResponse, error) {
	var familyItemsResponse FamilyItemsResponse
	if err := bggClient.get(ctx, "/family", fr, &familyItemsResponse); err != nil {
		return nil, err
	}
	return &familyItemsResponse, nil
}

func (bggClient *Client) GetForumList(ctx context.Context, fr *ForumListRequest) (*ForumListResponse, error) {
	var forumListResponse ForumListResponse
	if err := bggClient.get(ctx, "/forumlist", fr, &forumListResponse); err != nil {
		return nil, err
	}
	return &forumListResponse, nil
}

// GetForum returns an error matching ErrNotFound for an unknown forum id
func (bggClient *Client) GetForum(ctx context.Context, fr *ForumRequest) (*ForumResponse, error) {
	var forumResponse ForumResponse
	if err := bggClient.get(ctx, "/forum", fr, &forumResponse); err != nil {
		// BGG responds to an unknown forum id with a 200 html page rather than xml
		err = markNotFound(err, func(e *ErrorResponse) bool {
			return e.Response.StatusCode == http.StatusOK && strings.Contains(e.Response.Header.Get("Content-Type"), "html")
		})
		return nil, fmt.Errorf("forum %d: %w", fr.ID, err)
	}
	return &forumResponse, nil
}

// GetThread returns an error matching ErrNotFound for an unknown thread id
func (bggClient *Client) GetThread(ctx context.Context, tr *ThreadRequest) (*ThreadResponse, error) {
	var threadResponse ThreadResponse
	if err := bggClient.get(ctx, "/thread", tr, &threadResponse); err != nil {
		return nil, fmt.Errorf("thread %d: %w", tr.ID, err)
	}
	return &threadResponse, nil
}

func (bggClient *Client) GetHotItems(ctx context.Context, hr *HotItemsRequest) (*HotItemsResponse, error) {
	var hotItemsResponse HotItemsResponse
	if err := bggClient.get(ctx, "/hot", hr, &hotItemsResponse); err != nil {
		return nil, err
	}
	return &hotItemsResponse, nil
}

// GetUser returns an error matching ErrNotFound for an unknown username
func (bggClient *Client) GetUser(ctx context.Context, ur *UserRequest) (*UserResponse, error) {
	var userResponse UserResponse
	if err := bggClient.get(ctx, "/user", ur, &userResponse); err != nil {
		return nil, fmt.Errorf("user %s: %w", ur.Name, err)
	}
	return &userResponse, nil
}

// GetGuild returns an error matching ErrNotFound for an unknown guild id
func (bggClient *Client) GetGuild(ctx context.Context, gr *GuildRequest) (*GuildResponse, error) {
	var guildResponse GuildResponse
	if err := bggClient.get(ctx, "/guild", gr, &guildResponse); err != nil {
		return nil, fmt.Errorf("guild %d: %w", gr.ID, err)
	}
	return &guildResponse, nil
}

// GetCollection returns an *AcceptedError matching ErrCollectionNotReady while BGG is still building
// the collection; callers should retry the same request after a few seconds.
// An unknown username returns an error matching ErrNotFound.
func (bggClient *Client) GetCollection(ctx context.Context, cr *CollectionRequest) (*CollectionResponse, error) {
	var collectionResponse CollectionResponse
	if err := bggClient.get(ctx, "/collection", cr, &collectionResponse); err != nil {
		// BGG responds to an unknown username with a 200 <errors> document
		err = markNotFound(err, func(e *ErrorResponse) bool {
			return strings.Contains(e.Message, "Invalid username specified")
		})
		return nil, fmt.Errorf("collection for %s: %w", cr.Username, err)
	}
	return &collectionResponse, nil
}
