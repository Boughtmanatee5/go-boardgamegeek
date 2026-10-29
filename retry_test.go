package boardgamegeek

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newSequenceClient returns a Client whose nth request is answered with statuses[n]
// (the last status repeats) and a pointer to the number of requests made
func newSequenceClient(t *testing.T, maxRetries int, statuses ...int) (*Client, *int) {
	t.Helper()
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := statuses[min(requests, len(statuses)-1)]
		requests++
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`<items total="0"></items>`))
	}))
	t.Cleanup(server.Close)

	baseURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	c := NewClient(WithBaseURL(baseURL), WithMaxRetries(maxRetries))
	c.retryWait = time.Millisecond
	return c, &requests
}

func TestRetriesAcceptedUntilOK(t *testing.T) {
	c, requests := newSequenceClient(t, 3, http.StatusAccepted, http.StatusAccepted, http.StatusOK)

	_, err := c.Search(t.Context(), &SearchOptions{Query: "x"})
	assert.NoError(t, err)
	assert.Equal(t, 3, *requests)
}

func TestRetriesTooManyRequestsUntilOK(t *testing.T) {
	c, requests := newSequenceClient(t, 3, http.StatusTooManyRequests, http.StatusOK)

	_, err := c.Search(t.Context(), &SearchOptions{Query: "x"})
	assert.NoError(t, err)
	assert.Equal(t, 2, *requests)
}

func TestRetriesGiveUpAfterMaxRetries(t *testing.T) {
	c, requests := newSequenceClient(t, 2, http.StatusTooManyRequests)

	_, err := c.Search(t.Context(), &SearchOptions{Query: "x"})
	var rateLimitErr *RateLimitError
	assert.ErrorAs(t, err, &rateLimitErr)
	assert.Equal(t, 3, *requests, "one request plus two retries")
}

func TestNoRetriesByDefault(t *testing.T) {
	c, requests := newSequenceClient(t, 0, http.StatusAccepted, http.StatusOK)

	_, err := c.Search(t.Context(), &SearchOptions{Query: "x"})
	var acceptedErr *AcceptedError
	assert.ErrorAs(t, err, &acceptedErr)
	assert.Equal(t, 1, *requests)
}

func TestOtherErrorsAreNotRetried(t *testing.T) {
	c, requests := newSequenceClient(t, 3, http.StatusInternalServerError, http.StatusOK)

	_, err := c.Search(t.Context(), &SearchOptions{Query: "x"})
	assert.Error(t, err)
	assert.Equal(t, 1, *requests)
}

func TestRetryWaitStopsWhenContextIsCancelled(t *testing.T) {
	c, requests := newSequenceClient(t, 3, http.StatusAccepted)
	c.retryWait = time.Hour

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	_, err := c.Search(ctx, &SearchOptions{Query: "x"})
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, 1, *requests)
}

func TestRetryDelay(t *testing.T) {
	c := NewClient(WithMaxRetries(10))
	resp := &http.Response{}

	wait, retry := c.retryDelay(&RateLimitError{Response: resp, RetryAfter: 30 * time.Second}, 0)
	assert.True(t, retry)
	assert.Equal(t, 30*time.Second, wait, "Retry-After wins over backoff")

	wait, _ = c.retryDelay(&AcceptedError{Response: resp}, 2)
	assert.Equal(t, 8*time.Second, wait, "2s doubled twice")

	wait, _ = c.retryDelay(&AcceptedError{Response: resp}, 9)
	assert.Equal(t, maxRetryWait, wait)

	_, retry = c.retryDelay(&AcceptedError{Response: resp}, 10)
	assert.False(t, retry)
}
