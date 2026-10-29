package integration_test

import (
	"testing"

	boardgamegeek "github.com/Boughtmanatee5/BoardGameGeekClient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BGG's hot lists always have 50 items
const hotListSize = 50

func TestGetHotItemsDefaultsToBoardGames(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetHotItems(ctx, &boardgamegeek.HotItemsRequest{})
	require.NoError(t, err)
	assert.Len(t, response.Items, hotListSize)
	assert.NotEmpty(t, response.TermsOfUse)

	for i, item := range response.Items {
		assert.Equal(t, i+1, item.Rank)
	}
}

func TestGetHotItemsPopulatesFields(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.HotItemsRequest{
		Type: boardgamegeek.HotItemTypeBoardGame,
	}

	response, err := client.GetHotItems(ctx, request)
	require.NoError(t, err)
	require.NotEmpty(t, response.Items)

	item := response.Items[0]
	assert.NotZero(t, item.Id)
	assert.NotEmpty(t, item.Name.Value)
	assert.NotEmpty(t, item.Thumbnail.Value)
	require.NotNil(t, item.YearPublished)
	assert.NotEmpty(t, item.YearPublished.Value)
}

func TestGetHotItemsByRPGType(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	boardGames, err := client.GetHotItems(ctx, &boardgamegeek.HotItemsRequest{Type: boardgamegeek.HotItemTypeBoardGame})
	require.NoError(t, err)
	rpgs, err := client.GetHotItems(ctx, &boardgamegeek.HotItemsRequest{Type: boardgamegeek.HotItemTypeRPG})
	require.NoError(t, err)

	assert.Len(t, rpgs.Items, hotListSize)
	assert.NotEqual(t, boardGames.Items[0].Id, rpgs.Items[0].Id)
}

func TestGetHotItemsByPersonTypeHasNoYearPublished(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.HotItemsRequest{
		Type: boardgamegeek.HotItemTypeBoardGamePerson,
	}

	response, err := client.GetHotItems(ctx, request)
	require.NoError(t, err)
	require.NotEmpty(t, response.Items)
	for _, item := range response.Items {
		assert.NotEmpty(t, item.Name.Value)
		assert.Nil(t, item.YearPublished)
	}
}
