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
	var threadErr threadError
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
