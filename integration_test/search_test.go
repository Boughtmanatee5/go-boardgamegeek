//go:build integration

package integration_test

import (
	"testing"

	"github.com/Boughtmanatee5/go-boardgamegeek"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchOptions{
		Query: "guards of atlantis II",
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, int(response.Total), len(response.Items))
}

func TestSearchByType(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchOptions{
		Query: "guards of atlantis II",
		Types: []boardgamegeek.SearchType{boardgamegeek.SearchTypeBoardGame},
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, int(response.Total), len(response.Items))

	for _, item := range response.Items {
		assert.Equal(t, string(boardgamegeek.SearchTypeBoardGame), item.Type)
	}
}

func TestSearchByExact(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchOptions{
		Query: "guards of atlantis II",
		Exact: true,
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Items))
	// gaurds of atlantis 2 id
	assert.Equal(t, int64(267609), response.Items[0].ID)
	require.NotNil(t, response.Items[0].YearPublished)
	assert.EqualValues(t, 2022, *response.Items[0].YearPublished)
}

func TestSearchByMultipleTypes(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchOptions{
		Query: "catan",
		Types: []boardgamegeek.SearchType{boardgamegeek.SearchTypeBoardGame, boardgamegeek.SearchTypeBoardGameExpansion},
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, int(response.Total), len(response.Items))

	for _, item := range response.Items {
		assert.Contains(t, []string{string(boardgamegeek.SearchTypeBoardGame), string(boardgamegeek.SearchTypeBoardGameExpansion)}, item.Type)
	}
}

func TestSearchByDesignerType(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchOptions{
		Query: "reiner knizia",
		Types: []boardgamegeek.SearchType{boardgamegeek.SearchTypeBoardGameDesigner},
		Exact: true,
	}

	response, err := client.Search(ctx, searchRequest)
	require.NoError(t, err)
	require.Len(t, response.Items, 1)
	// designers have their own ids: 2 is Reiner Knizia, while thing 2 is the game Dragonmaster
	assert.Equal(t, int64(2), response.Items[0].ID)
	assert.Equal(t, string(boardgamegeek.SearchTypeBoardGameDesigner), response.Items[0].Type)
	require.NotNil(t, response.Items[0].Name)
	assert.Equal(t, "Reiner Knizia", response.Items[0].Name.Value)
}

func TestSearchByRPGIssueType(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchOptions{
		Query: "dragon",
		Types: []boardgamegeek.SearchType{boardgamegeek.SearchTypeRPGIssue},
	}

	response, err := client.Search(ctx, searchRequest)
	require.NoError(t, err)
	assert.NotEmpty(t, response.Items)

	for _, item := range response.Items {
		assert.Equal(t, string(boardgamegeek.SearchTypeRPGIssue), item.Type)
	}
}

// Edge cases

func TestSearchWithNoResults(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchOptions{
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
	searchRequest := &boardgamegeek.SearchOptions{
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

	lowercaseResponse, err := client.Search(ctx, &boardgamegeek.SearchOptions{Query: "catan"})
	assert.NoError(t, err)

	uppercaseResponse, err := client.Search(ctx, &boardgamegeek.SearchOptions{Query: "CATAN"})
	assert.NoError(t, err)

	assert.Equal(t, lowercaseResponse.Total, uppercaseResponse.Total)
}

func TestSearchWithSpecialCharactersInQuery(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchOptions{
		Query: "Sid Meier's Civilization",
		Exact: true,
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Items))
	assert.Equal(t, int64(69188), response.Items[0].ID)
}

func TestSearchByExactWithNoMatch(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	searchRequest := &boardgamegeek.SearchOptions{
		// A partial, non-exact name for Catan; exact matching should find nothing.
		Query: "catan the game that does not exist",
		Exact: true,
	}

	response, err := client.Search(ctx, searchRequest)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), response.Total)
	assert.Empty(t, response.Items)
}
