package integration_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	boardgamegeek "github.com/Boughtmanatee5/BoardGameGeekClient"
	"golang.org/x/oauth2"
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

// testThrottle spaces out tests so the suite stays under BGG's rate limit,
// which returns 429 Too Many Requests after roughly 40 requests in a burst.
const testThrottle = 300 * time.Millisecond

func newTestClient() *boardgamegeek.Client {
	time.Sleep(testThrottle)

	token := &oauth2.Token{
		TokenType:   "Bearer",
		AccessToken: bggToken,
	}
	return boardgamegeek.NewClient(token)
}
