package boardgamegeek

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

const boadgamegeekURL = "https://boardgamegeek.com/xmlapi2"

// apiError is returned with HTTP 200 in place of the expected document when an entity
// does not exist, e.g. <error message='Thread Not Found' /> from /thread and /guild
type apiError struct {
	XMLName xml.Name `xml:"error"`
	Message string   `xml:"message,attr"`
}

type Client struct {
	httpClient http.Client
	token      *oauth2.Token
}

func NewClient(token *oauth2.Token) *Client {
	httpClient := http.Client{}

	return &Client{
		httpClient: httpClient,
		token:      token,
	}
}

func (bggClient *Client) Search(ctx context.Context, sr *SearchRequest) (*SearchResponse, error) {
	searchURL, urlParseErr := url.Parse(boadgamegeekURL)
	if urlParseErr != nil {
		return nil, fmt.Errorf("error parsing boardgamegeek api host %w", urlParseErr)
	}
	searchURL = searchURL.JoinPath("/search")
	searchURL.RawQuery = sr.toQueryString()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, searchURL.String(), nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error creating search request %w", reqErr)
	}

	if bggClient.token != nil {
		bggClient.token.SetAuthHeader(req)
	}

	resp, bggErr := bggClient.httpClient.Do(req)
	if bggErr != nil {
		return nil, fmt.Errorf("recieved error response from boardgamegeek %w", bggErr)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("error reading response body %w", readErr)
	}
	var searchResponse *SearchResponse
	parseErr := xml.Unmarshal(body, &searchResponse)
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing xml body %w", parseErr)
	}

	return searchResponse, nil
}

func (bggClient *Client) GetThings(ctx context.Context, gr *GetThingsRequest) (*GetThingsResponse, error) {
	getURL, urlParseErr := url.Parse(boadgamegeekURL)
	if urlParseErr != nil {
		return nil, fmt.Errorf("error parsing boardgamegeek url %w", urlParseErr)
	}

	getURL = getURL.JoinPath("/thing")
	getURL.RawQuery = gr.toQueryString()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, getURL.String(), nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error creating get things request %w", reqErr)
	}

	if bggClient.token != nil {
		bggClient.token.SetAuthHeader(req)
	}

	resp, bggErr := bggClient.httpClient.Do(req)
	if bggErr != nil {
		return nil, fmt.Errorf("recieved error response from boardgamegeek %w", bggErr)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("error reading response body %w", readErr)
	}
	var GetThingsResponse *GetThingsResponse
	parseErr := xml.Unmarshal(body, &GetThingsResponse)
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing xml body %w", parseErr)
	}

	return GetThingsResponse, nil
}

func (bggClient *Client) GetFamilyItems(ctx context.Context, fr *FamilyItemRequest) (*FamilyItemsResponse, error) {
	familyURL, urlParseErr := url.Parse(boadgamegeekURL)
	if urlParseErr != nil {
		return nil, fmt.Errorf("error parsing boardgamegeek url %w", urlParseErr)
	}

	familyURL = familyURL.JoinPath("/family")
	familyURL.RawQuery = fr.toQueryString()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, familyURL.String(), nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error creating get family items request %w", reqErr)
	}

	if bggClient.token != nil {
		bggClient.token.SetAuthHeader(req)
	}

	resp, bggErr := bggClient.httpClient.Do(req)
	if bggErr != nil {
		return nil, fmt.Errorf("recieved error response from boardgamegeek %w", bggErr)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("error reading response body %w", readErr)
	}
	var familyItemsResponse *FamilyItemsResponse
	parseErr := xml.Unmarshal(body, &familyItemsResponse)
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing xml body %w", parseErr)
	}

	return familyItemsResponse, nil
}

func (bggClient *Client) GetForumList(ctx context.Context, fr *ForumListRequest) (*ForumListResponse, error) {
	forumListURL, urlParseErr := url.Parse(boadgamegeekURL)
	if urlParseErr != nil {
		return nil, fmt.Errorf("error parsing boardgamegeek url %w", urlParseErr)
	}

	forumListURL = forumListURL.JoinPath("/forumlist")
	forumListURL.RawQuery = fr.toQueryString()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, forumListURL.String(), nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error creating get forum list request %w", reqErr)
	}

	if bggClient.token != nil {
		bggClient.token.SetAuthHeader(req)
	}

	resp, bggErr := bggClient.httpClient.Do(req)
	if bggErr != nil {
		return nil, fmt.Errorf("recieved error response from boardgamegeek %w", bggErr)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("error reading response body %w", readErr)
	}
	var forumListResponse *ForumListResponse
	parseErr := xml.Unmarshal(body, &forumListResponse)
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing xml body %w", parseErr)
	}

	return forumListResponse, nil
}

func (bggClient *Client) GetForum(ctx context.Context, fr *ForumRequest) (*ForumResponse, error) {
	forumURL, urlParseErr := url.Parse(boadgamegeekURL)
	if urlParseErr != nil {
		return nil, fmt.Errorf("error parsing boardgamegeek url %w", urlParseErr)
	}

	forumURL = forumURL.JoinPath("/forum")
	forumURL.RawQuery = fr.toQueryString()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, forumURL.String(), nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error creating get forum request %w", reqErr)
	}

	if bggClient.token != nil {
		bggClient.token.SetAuthHeader(req)
	}

	resp, bggErr := bggClient.httpClient.Do(req)
	if bggErr != nil {
		return nil, fmt.Errorf("recieved error response from boardgamegeek %w", bggErr)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	// BGG responds to an unknown forum id with a 200 html page rather than xml
	if !strings.Contains(resp.Header.Get("Content-Type"), "xml") {
		return nil, fmt.Errorf("forum %d not found", fr.ID)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("error reading response body %w", readErr)
	}
	var forumResponse *ForumResponse
	parseErr := xml.Unmarshal(body, &forumResponse)
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing xml body %w", parseErr)
	}

	return forumResponse, nil
}

func (bggClient *Client) GetThread(ctx context.Context, tr *ThreadRequest) (*ThreadResponse, error) {
	threadURL, urlParseErr := url.Parse(boadgamegeekURL)
	if urlParseErr != nil {
		return nil, fmt.Errorf("error parsing boardgamegeek url %w", urlParseErr)
	}

	threadURL = threadURL.JoinPath("/thread")
	threadURL.RawQuery = tr.toQueryString()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, threadURL.String(), nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error creating get thread request %w", reqErr)
	}

	if bggClient.token != nil {
		bggClient.token.SetAuthHeader(req)
	}

	resp, bggErr := bggClient.httpClient.Do(req)
	if bggErr != nil {
		return nil, fmt.Errorf("recieved error response from boardgamegeek %w", bggErr)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("error reading response body %w", readErr)
	}

	// BGG responds to an unknown thread id with a 200 <error> document
	var threadErr apiError
	if xml.Unmarshal(body, &threadErr) == nil {
		return nil, fmt.Errorf("thread %d: %s", tr.ID, threadErr.Message)
	}

	var threadResponse *ThreadResponse
	parseErr := xml.Unmarshal(body, &threadResponse)
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing xml body %w", parseErr)
	}

	return threadResponse, nil
}

func (bggClient *Client) GetHotItems(ctx context.Context, hr *HotItemsRequest) (*HotItemsResponse, error) {
	hotURL, urlParseErr := url.Parse(boadgamegeekURL)
	if urlParseErr != nil {
		return nil, fmt.Errorf("error parsing boardgamegeek url %w", urlParseErr)
	}

	hotURL = hotURL.JoinPath("/hot")
	hotURL.RawQuery = hr.toQueryString()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, hotURL.String(), nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error creating get hot items request %w", reqErr)
	}

	if bggClient.token != nil {
		bggClient.token.SetAuthHeader(req)
	}

	resp, bggErr := bggClient.httpClient.Do(req)
	if bggErr != nil {
		return nil, fmt.Errorf("recieved error response from boardgamegeek %w", bggErr)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("error reading response body %w", readErr)
	}
	var hotItemsResponse *HotItemsResponse
	parseErr := xml.Unmarshal(body, &hotItemsResponse)
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing xml body %w", parseErr)
	}

	return hotItemsResponse, nil
}

func (bggClient *Client) GetUser(ctx context.Context, ur *UserRequest) (*UserResponse, error) {
	userURL, urlParseErr := url.Parse(boadgamegeekURL)
	if urlParseErr != nil {
		return nil, fmt.Errorf("error parsing boardgamegeek url %w", urlParseErr)
	}

	userURL = userURL.JoinPath("/user")
	userURL.RawQuery = ur.toQueryString()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, userURL.String(), nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error creating get user request %w", reqErr)
	}

	if bggClient.token != nil {
		bggClient.token.SetAuthHeader(req)
	}

	resp, bggErr := bggClient.httpClient.Do(req)
	if bggErr != nil {
		return nil, fmt.Errorf("recieved error response from boardgamegeek %w", bggErr)
	}
	defer func() { _ = resp.Body.Close() }()

	// BGG responds to an unknown username with a 404
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("error reading response body %w", readErr)
	}
	var userResponse *UserResponse
	parseErr := xml.Unmarshal(body, &userResponse)
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing xml body %w", parseErr)
	}

	return userResponse, nil
}

func (bggClient *Client) GetGuild(ctx context.Context, gr *GuildRequest) (*GuildResponse, error) {
	guildURL, urlParseErr := url.Parse(boadgamegeekURL)
	if urlParseErr != nil {
		return nil, fmt.Errorf("error parsing boardgamegeek url %w", urlParseErr)
	}

	guildURL = guildURL.JoinPath("/guild")
	guildURL.RawQuery = gr.toQueryString()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, guildURL.String(), nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error creating get guild request %w", reqErr)
	}

	if bggClient.token != nil {
		bggClient.token.SetAuthHeader(req)
	}

	resp, bggErr := bggClient.httpClient.Do(req)
	if bggErr != nil {
		return nil, fmt.Errorf("recieved error response from boardgamegeek %w", bggErr)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("error reading response body %w", readErr)
	}

	// BGG responds to an unknown guild id with a 200 <error> document
	var guildErr apiError
	if xml.Unmarshal(body, &guildErr) == nil {
		return nil, fmt.Errorf("guild %d: %s", gr.ID, guildErr.Message)
	}

	var guildResponse *GuildResponse
	parseErr := xml.Unmarshal(body, &guildResponse)
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing xml body %w", parseErr)
	}

	return guildResponse, nil
}

// GetCollection returns ErrCollectionNotReady while BGG is still building the collection;
// callers should retry the same request after a few seconds.
func (bggClient *Client) GetCollection(ctx context.Context, cr *CollectionRequest) (*CollectionResponse, error) {
	collectionURL, urlParseErr := url.Parse(boadgamegeekURL)
	if urlParseErr != nil {
		return nil, fmt.Errorf("error parsing boardgamegeek url %w", urlParseErr)
	}

	collectionURL = collectionURL.JoinPath("/collection")
	collectionURL.RawQuery = cr.toQueryString()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, collectionURL.String(), nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error creating get collection request %w", reqErr)
	}

	if bggClient.token != nil {
		bggClient.token.SetAuthHeader(req)
	}

	resp, bggErr := bggClient.httpClient.Do(req)
	if bggErr != nil {
		return nil, fmt.Errorf("recieved error response from boardgamegeek %w", bggErr)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusAccepted {
		return nil, ErrCollectionNotReady
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("error reading response body %w", readErr)
	}

	// BGG responds to an invalid request, e.g. an unknown username, with a 200 <errors> document
	var collectionErrs collectionErrors
	if xml.Unmarshal(body, &collectionErrs) == nil {
		return nil, fmt.Errorf("collection for %s: %s", cr.Username, strings.Join(collectionErrs.Messages, "; "))
	}

	var collectionResponse *CollectionResponse
	parseErr := xml.Unmarshal(body, &collectionResponse)
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing xml body %w", parseErr)
	}

	return collectionResponse, nil
}
