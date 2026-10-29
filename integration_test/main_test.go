//go:build integration

package integration_test

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Boughtmanatee5/go-boardgamegeek"
)

const tokenEnvVar = "BGG_API_TOKEN"

var bggToken string

func TestMain(m *testing.M) {
	bggToken = os.Getenv(tokenEnvVar)
	if bggToken == "" {
		fmt.Fprintf(os.Stderr, "%s must be set to run integration tests\n", tokenEnvVar)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// requestInterval spaces out requests so the suite stays under BGG's rate limit,
// which returns 429 Too Many Requests after roughly 40 requests in a burst and refills slowly.
const requestInterval = 2 * time.Second

// collectionInterval spaces out collection requests: /collection has a tighter rate limit
// than the other endpoints and returns 429 after about 6 requests 2s apart.
const collectionInterval = 5 * time.Second

// throttledTransport waits until at least requestInterval (collectionInterval for /collection)
// has passed since the previous request before sending the next one
type throttledTransport struct {
	mu   sync.Mutex
	last time.Time
}

func (t *throttledTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	interval := requestInterval
	if strings.HasSuffix(req.URL.Path, "/collection") {
		interval = collectionInterval
	}

	t.mu.Lock()
	time.Sleep(time.Until(t.last.Add(interval)))
	t.last = time.Now()
	t.mu.Unlock()

	return http.DefaultTransport.RoundTrip(req)
}

// testHTTPClient is shared by every test client so the throttle applies across the whole suite
var testHTTPClient = &http.Client{Transport: &throttledTransport{}}

// newTestClient returns an authenticated, throttled client with any extra options applied
func newTestClient(options ...boardgamegeek.ClientOption) *boardgamegeek.Client {
	return boardgamegeek.NewClient(append([]boardgamegeek.ClientOption{
		boardgamegeek.WithHTTPClient(testHTTPClient),
		boardgamegeek.WithAuthToken(bggToken),
	}, options...)...)
}
