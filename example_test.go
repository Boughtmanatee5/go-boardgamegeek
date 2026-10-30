package boardgamegeek_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Boughtmanatee5/go-boardgamegeek"
)

// These examples call the live API, so they have no Output comment and are compiled but not run by go test.

func ExampleNewClient() {
	client := boardgamegeek.NewClient(
		boardgamegeek.WithAuthToken(os.Getenv("BGG_API_TOKEN")),
		boardgamegeek.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
		boardgamegeek.WithMaxRetries(5),
	)

	resp, err := client.Search(context.Background(), &boardgamegeek.SearchOptions{
		Query: "Catan",
		Types: []boardgamegeek.SearchType{boardgamegeek.SearchTypeBoardGame},
		Exact: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range resp.Items {
		fmt.Println(item.ID, item.Name.Value)
	}
}

func ExampleClient_GetThings() {
	client := boardgamegeek.NewClient(boardgamegeek.WithAuthToken(os.Getenv("BGG_API_TOKEN")))

	resp, err := client.GetThings(context.Background(), &boardgamegeek.ThingsOptions{
		IDs:   []int64{13, 822},
		Stats: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, thing := range resp.Items {
		fmt.Println(thing.ID, thing.Names[0].Value)
	}
}

func ExampleClient_GetCollection() {
	// BGG answers 202 while it builds a collection; WithMaxRetries waits and retries for you
	client := boardgamegeek.NewClient(
		boardgamegeek.WithAuthToken(os.Getenv("BGG_API_TOKEN")),
		boardgamegeek.WithMaxRetries(10),
	)

	resp, err := client.GetCollection(context.Background(), &boardgamegeek.CollectionOptions{
		Username: "Aldie",
		Own:      new(true),
	})
	if errors.Is(err, boardgamegeek.ErrNotFound) {
		log.Fatal("no such user")
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.TotalItems, "games owned")
}

func ExampleClient_AllForumThreads() {
	client := boardgamegeek.NewClient(
		boardgamegeek.WithAuthToken(os.Getenv("BGG_API_TOKEN")),
		boardgamegeek.WithMaxRetries(5),
	)

	// Pages are requested as the loop needs them; breaking out stops further requests
	for thread, err := range client.AllForumThreads(context.Background(), &boardgamegeek.ForumOptions{ID: 1508}) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(thread.Subject)
	}
}

func ExampleErrNotFound() {
	client := boardgamegeek.NewClient(boardgamegeek.WithAuthToken(os.Getenv("BGG_API_TOKEN")))

	_, err := client.GetThread(context.Background(), &boardgamegeek.ThreadOptions{ID: 999999999})

	var rateLimitErr *boardgamegeek.RateLimitError
	switch {
	case errors.Is(err, boardgamegeek.ErrNotFound):
		fmt.Println("thread does not exist")
	case errors.As(err, &rateLimitErr):
		fmt.Println("rate limited, retry after", rateLimitErr.RetryAfter)
	case err != nil:
		log.Fatal(err)
	}
}
