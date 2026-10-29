package integration_test

import (
	"testing"

	boardgamegeek "github.com/Boughtmanatee5/BoardGameGeekClient"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchRequest{
		Query: "guards of atlantis II",
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, int(response.Total), len(response.Items))
}

func TestSearchByThingType(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchRequest{
		Query: "guards of atlantis II",
		Types: []boardgamegeek.ThingType{boardgamegeek.BOARDGAME},
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, int(response.Total), len(response.Items))

	for _, item := range response.Items {
		assert.Equal(t, string(boardgamegeek.BOARDGAME), item.Type)
	}
}

func TestSearchByExact(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchRequest{
		Query: "guards of atlantis II",
		Exact: true,
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Items))
	// gaurds of atlantis 2 id
	assert.Equal(t, int64(267609), response.Items[0].Id)
}

func TestSearchByMultipleTypes(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchRequest{
		Query: "catan",
		Types: []boardgamegeek.ThingType{boardgamegeek.BOARDGAME, boardgamegeek.BOARDGAMEEXPANSION},
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, int(response.Total), len(response.Items))

	for _, item := range response.Items {
		assert.Contains(t, []string{string(boardgamegeek.BOARDGAME), string(boardgamegeek.BOARDGAMEEXPANSION)}, item.Type)
	}
}

// Edge cases

func TestSearchWithNoResults(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchRequest{
		Query: "zzzzznonexistentgamexyz123",
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), response.Total)
	assert.Empty(t, response.Items)
}

func TestSearchWithEmptyQuery(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchRequest{
		Query: "",
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), response.Total)
	assert.Empty(t, response.Items)
}

func TestSearchIsCaseInsensitive(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	lowercaseResponse, err := client.Search(ctx, &boardgamegeek.SearchRequest{Query: "catan"})
	assert.NoError(t, err)

	uppercaseResponse, err := client.Search(ctx, &boardgamegeek.SearchRequest{Query: "CATAN"})
	assert.NoError(t, err)

	assert.Equal(t, lowercaseResponse.Total, uppercaseResponse.Total)
}

func TestSearchWithSpecialCharactersInQuery(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchRequest{
		Query: "Sid Meier's Civilization",
		Exact: true,
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Items))
	assert.Equal(t, int64(69188), response.Items[0].Id)
}

func TestSearchByExactWithNoMatch(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchRequest{
		// A partial, non-exact name for Catan; exact matching should find nothing.
		Query: "catan the game that does not exist",
		Exact: true,
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), response.Total)
	assert.Empty(t, response.Items)
}
