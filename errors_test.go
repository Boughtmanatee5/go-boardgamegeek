package boardgamegeek_test

import (
	"errors"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Boughtmanatee5/go-boardgamegeek"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests cover error paths that BGG can't be made to produce on demand, such as 429.

// newStubClient returns a Client whose requests are answered with status, header and body
func newStubClient(t *testing.T, status int, header http.Header, body string) *boardgamegeek.Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		maps.Copy(w.Header(), header)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	baseURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	return boardgamegeek.NewClient(boardgamegeek.WithBaseURL(baseURL))
}

func xmlHeader() http.Header {
	return http.Header{"Content-Type": {"text/xml; charset=utf-8"}}
}

func TestAcceptedMatchesErrCollectionNotReady(t *testing.T) {
	client := newStubClient(t, http.StatusAccepted, xmlHeader(), "")

	_, err := client.GetCollection(t.Context(), &boardgamegeek.CollectionOptions{Username: "someone"})
	var accepted *boardgamegeek.AcceptedError
	assert.ErrorAs(t, err, &accepted)
	assert.ErrorIs(t, err, boardgamegeek.ErrCollectionNotReady)
}

func TestTooManyRequestsReturnsRateLimitError(t *testing.T) {
	client := newStubClient(t, http.StatusTooManyRequests, http.Header{"Retry-After": {"30"}}, "")

	_, err := client.GetThings(t.Context(), &boardgamegeek.ThingsOptions{IDs: []int64{1}})
	var rateLimitErr *boardgamegeek.RateLimitError
	require.ErrorAs(t, err, &rateLimitErr)
	assert.Equal(t, 30*time.Second, rateLimitErr.RetryAfter)
	assert.False(t, errors.Is(err, boardgamegeek.ErrNotFound))
}

func TestServerErrorIsNotErrNotFound(t *testing.T) {
	client := newStubClient(t, http.StatusInternalServerError, nil, "")

	_, err := client.GetThings(t.Context(), &boardgamegeek.ThingsOptions{IDs: []int64{1}})
	var errResp *boardgamegeek.ErrorResponse
	require.ErrorAs(t, err, &errResp)
	assert.Equal(t, http.StatusInternalServerError, errResp.Response.StatusCode)
	assert.False(t, errors.Is(err, boardgamegeek.ErrNotFound))
}

func TestErrorDocumentMatchesErrNotFound(t *testing.T) {
	client := newStubClient(t, http.StatusOK, xmlHeader(), `<error message="Guild not found" />`)

	_, err := client.GetGuild(t.Context(), &boardgamegeek.GuildOptions{ID: 1})
	var errResp *boardgamegeek.ErrorResponse
	require.ErrorAs(t, err, &errResp)
	assert.Equal(t, "Guild not found", errResp.Message)
	assert.ErrorIs(t, err, boardgamegeek.ErrNotFound)
}

func TestCollectionOtherErrorsAreNotErrNotFound(t *testing.T) {
	client := newStubClient(t, http.StatusOK, xmlHeader(), `<errors><error><message>Invalid parameter</message></error></errors>`)

	_, err := client.GetCollection(t.Context(), &boardgamegeek.CollectionOptions{Username: "someone"})
	var errResp *boardgamegeek.ErrorResponse
	require.ErrorAs(t, err, &errResp)
	assert.Equal(t, "Invalid parameter", errResp.Message)
	assert.False(t, errors.Is(err, boardgamegeek.ErrNotFound))
}

func TestHTMLIsErrNotFoundOnlyForForum(t *testing.T) {
	htmlHeader := http.Header{"Content-Type": {"text/html; charset=utf-8"}}

	forumClient := newStubClient(t, http.StatusOK, htmlHeader, "<html></html>")
	_, err := forumClient.GetForum(t.Context(), &boardgamegeek.ForumOptions{ID: 1})
	assert.ErrorIs(t, err, boardgamegeek.ErrNotFound)

	searchClient := newStubClient(t, http.StatusOK, htmlHeader, "<html></html>")
	_, err = searchClient.Search(t.Context(), &boardgamegeek.SearchOptions{Query: "x"})
	var errResp *boardgamegeek.ErrorResponse
	assert.ErrorAs(t, err, &errResp)
	assert.False(t, errors.Is(err, boardgamegeek.ErrNotFound))
}

func TestPlaysMessageBoxMatchesErrNotFound(t *testing.T) {
	body := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\t<div class='messagebox error'>\n\t\tInvalid object or user\n\t</div>"
	client := newStubClient(t, http.StatusOK, xmlHeader(), body)

	_, err := client.GetPlays(t.Context(), &boardgamegeek.PlaysOptions{Username: "nobody"})
	var errResp *boardgamegeek.ErrorResponse
	require.ErrorAs(t, err, &errResp)
	assert.Equal(t, "Invalid object or user", errResp.Message)
	assert.ErrorIs(t, err, boardgamegeek.ErrNotFound)
}

func TestMessageBoxIsErrNotFoundOnlyForPlays(t *testing.T) {
	client := newStubClient(t, http.StatusOK, xmlHeader(), "<div class='messagebox error'>Invalid object or user</div>")

	_, err := client.Search(t.Context(), &boardgamegeek.SearchOptions{Query: "x"})
	var errResp *boardgamegeek.ErrorResponse
	require.ErrorAs(t, err, &errResp)
	assert.False(t, errors.Is(err, boardgamegeek.ErrNotFound))
}

func TestPlaysWithoutUserOrIDFailsBeforeRequest(t *testing.T) {
	// BGG would answer <error message='Not Found'/>; GetPlays must not send the request or report not found
	client := newStubClient(t, http.StatusOK, xmlHeader(), `<error message="Not Found" />`)

	_, err := client.GetPlays(t.Context(), &boardgamegeek.PlaysOptions{Subtype: boardgamegeek.ThingTypeBoardGame})
	require.Error(t, err)
	var errResp *boardgamegeek.ErrorResponse
	assert.False(t, errors.As(err, &errResp))
	assert.False(t, errors.Is(err, boardgamegeek.ErrNotFound))
}
