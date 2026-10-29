package integration_test

import (
	"fmt"
	"os"
	"testing"

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

func newTestClient() *boardgamegeek.Client {
	token := &oauth2.Token{
		TokenType:   "Bearer",
		AccessToken: bggToken,
	}
	return boardgamegeek.NewClient(token)
}
