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
)

const (
	defaultBaseURL   = "https://boardgamegeek.com/xmlapi2/"
	defaultUserAgent = "BoardGameGeekClient-go"
)

// apiError is returned with HTTP 200 in place of the expected document when an entity
// does not exist, e.g. <error message='Thread Not Found' /> from /thread and /guild
type apiError struct {
	XMLName xml.Name `xml:"error"`
	Message string   `xml:"message,attr"`
}

func (e *apiError) Error() string {
	return e.Message
}

// errAccepted is returned by get for HTTP 202, which BGG uses while it queues a request
var errAccepted = errors.New("request accepted by boardgamegeek but not yet processed")

// errNotXML is returned by get when BGG responds with a non-xml document, e.g. the
// html page /forum returns for an unknown id
var errNotXML = errors.New("boardgamegeek returned a non-xml response")

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

// get requests path with the given query and unmarshals the xml body into v.
// BGG's "not found" responses come back as errNotXML, *apiError or *collectionErrors,
// and HTTP 202 as errAccepted, so each endpoint can report them in its own terms.
func (bggClient *Client) get(ctx context.Context, path string, query string, v any) error {
	reqURL := bggClient.baseURL.JoinPath(path)
	reqURL.RawQuery = query

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

	if resp.StatusCode == http.StatusAccepted {
		return errAccepted
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "xml") {
		return errNotXML
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return fmt.Errorf("error reading response body %w", readErr)
	}

	var bggErrResp apiError
	if xml.Unmarshal(body, &bggErrResp) == nil {
		return &bggErrResp
	}

	var collectionErrs collectionErrors
	if xml.Unmarshal(body, &collectionErrs) == nil {
		return &collectionErrs
	}

	if parseErr := xml.Unmarshal(body, v); parseErr != nil {
		return fmt.Errorf("error parsing xml body %w", parseErr)
	}

	return nil
}

func (bggClient *Client) Search(ctx context.Context, sr *SearchRequest) (*SearchResponse, error) {
	var searchResponse SearchResponse
	if err := bggClient.get(ctx, "/search", sr.toQueryString(), &searchResponse); err != nil {
		return nil, err
	}
	return &searchResponse, nil
}

func (bggClient *Client) GetThings(ctx context.Context, gr *GetThingsRequest) (*GetThingsResponse, error) {
	var getThingsResponse GetThingsResponse
	if err := bggClient.get(ctx, "/thing", gr.toQueryString(), &getThingsResponse); err != nil {
		return nil, err
	}
	return &getThingsResponse, nil
}

func (bggClient *Client) GetFamilyItems(ctx context.Context, fr *FamilyItemRequest) (*FamilyItemsResponse, error) {
	var familyItemsResponse FamilyItemsResponse
	if err := bggClient.get(ctx, "/family", fr.toQueryString(), &familyItemsResponse); err != nil {
		return nil, err
	}
	return &familyItemsResponse, nil
}

func (bggClient *Client) GetForumList(ctx context.Context, fr *ForumListRequest) (*ForumListResponse, error) {
	var forumListResponse ForumListResponse
	if err := bggClient.get(ctx, "/forumlist", fr.toQueryString(), &forumListResponse); err != nil {
		return nil, err
	}
	return &forumListResponse, nil
}

func (bggClient *Client) GetForum(ctx context.Context, fr *ForumRequest) (*ForumResponse, error) {
	var forumResponse ForumResponse
	err := bggClient.get(ctx, "/forum", fr.toQueryString(), &forumResponse)
	// BGG responds to an unknown forum id with a 200 html page rather than xml
	if errors.Is(err, errNotXML) {
		return nil, fmt.Errorf("forum %d not found", fr.ID)
	}
	if err != nil {
		return nil, err
	}
	return &forumResponse, nil
}

func (bggClient *Client) GetThread(ctx context.Context, tr *ThreadRequest) (*ThreadResponse, error) {
	var threadResponse ThreadResponse
	err := bggClient.get(ctx, "/thread", tr.toQueryString(), &threadResponse)
	// BGG responds to an unknown thread id with a 200 <error> document
	var threadErr *apiError
	if errors.As(err, &threadErr) {
		return nil, fmt.Errorf("thread %d: %s", tr.ID, threadErr.Message)
	}
	if err != nil {
		return nil, err
	}
	return &threadResponse, nil
}

func (bggClient *Client) GetHotItems(ctx context.Context, hr *HotItemsRequest) (*HotItemsResponse, error) {
	var hotItemsResponse HotItemsResponse
	if err := bggClient.get(ctx, "/hot", hr.toQueryString(), &hotItemsResponse); err != nil {
		return nil, err
	}
	return &hotItemsResponse, nil
}

// GetUser returns a "404 Not Found" status error for an unknown username
func (bggClient *Client) GetUser(ctx context.Context, ur *UserRequest) (*UserResponse, error) {
	var userResponse UserResponse
	if err := bggClient.get(ctx, "/user", ur.toQueryString(), &userResponse); err != nil {
		return nil, err
	}
	return &userResponse, nil
}

func (bggClient *Client) GetGuild(ctx context.Context, gr *GuildRequest) (*GuildResponse, error) {
	var guildResponse GuildResponse
	err := bggClient.get(ctx, "/guild", gr.toQueryString(), &guildResponse)
	// BGG responds to an unknown guild id with a 200 <error> document
	var guildErr *apiError
	if errors.As(err, &guildErr) {
		return nil, fmt.Errorf("guild %d: %s", gr.ID, guildErr.Message)
	}
	if err != nil {
		return nil, err
	}
	return &guildResponse, nil
}

// GetCollection returns ErrCollectionNotReady while BGG is still building the collection;
// callers should retry the same request after a few seconds.
func (bggClient *Client) GetCollection(ctx context.Context, cr *CollectionRequest) (*CollectionResponse, error) {
	var collectionResponse CollectionResponse
	err := bggClient.get(ctx, "/collection", cr.toQueryString(), &collectionResponse)
	if errors.Is(err, errAccepted) {
		return nil, ErrCollectionNotReady
	}
	// BGG responds to an invalid request, e.g. an unknown username, with a 200 <errors> document
	var collectionErrs *collectionErrors
	if errors.As(err, &collectionErrs) {
		return nil, fmt.Errorf("collection for %s: %s", cr.Username, collectionErrs.Error())
	}
	if err != nil {
		return nil, err
	}
	return &collectionResponse, nil
}
