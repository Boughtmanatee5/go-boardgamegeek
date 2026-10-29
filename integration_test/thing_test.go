package integration_test

import (
	"testing"

	boardgamegeek "github.com/Boughtmanatee5/BoardGameGeekClient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Catan (base game)
const catanId int64 = 13

// Catan: Seafarers (expansion)
const catanSeafarersId int64 = 325

func TestGetThingsById(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs: []int64{catanId},
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 1)
	assert.Equal(t, catanId, response.Items[0].Id)
	assert.Equal(t, string(boardgamegeek.BOARDGAME), response.Items[0].Type)
}

func TestGetThingsByMultipleIds(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs: []int64{catanId, catanSeafarersId},
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 2)

	var returnedIds []int64
	for _, item := range response.Items {
		returnedIds = append(returnedIds, item.Id)
	}
	assert.ElementsMatch(t, []int64{catanId, catanSeafarersId}, returnedIds)
}

func TestGetThingsByType(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs:   []int64{catanId, catanSeafarersId},
		Types: []boardgamegeek.ThingType{boardgamegeek.BOARDGAMEEXPANSION},
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 1)
	assert.Equal(t, catanSeafarersId, response.Items[0].Id)
	assert.Equal(t, string(boardgamegeek.BOARDGAMEEXPANSION), response.Items[0].Type)
}

func TestGetThingsWithVersions(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs:      []int64{catanId},
		Versions: true,
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response.Items[0].Versions)
	assert.NotEmpty(t, response.Items[0].Versions.Items)
}

func TestGetThingsWithVideos(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs:    []int64{catanId},
		Videos: true,
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response.Items[0].Videos)
	assert.Greater(t, response.Items[0].Videos.Total, int64(0))
}

func TestGetThingsWithStats(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs:   []int64{catanId},
		Stats: true,
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response.Items[0].Statistics)
}

func TestGetThingsWithMarketPlace(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs:         []int64{catanId},
		MarketPlace: true,
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response.Items[0].MarketplaceListings)
}

func TestGetThingsWithComments(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs:      []int64{catanId},
		Comments: true,
		Page:     1,
		PageSize: 10,
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response.Items[0].Comments)
	assert.Len(t, response.Items[0].Comments.Comments, 10)
	assert.Greater(t, response.Items[0].Comments.TotalItems, int64(0))
}

func TestGetThingsWithRatingComments(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs:            []int64{catanId},
		RatingComments: true,
		Page:           1,
		PageSize:       10,
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response.Items[0].Comments)
	assert.Len(t, response.Items[0].Comments.Comments, 10)
	for _, comment := range response.Items[0].Comments.Comments {
		assert.NotEmpty(t, comment.Rating)
	}
}

func TestGetThingsByPage(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	firstPageRequest := &boardgamegeek.GetThingsRequest{
		IDs:      []int64{catanId},
		Comments: true,
		Page:     1,
		PageSize: 10,
	}
	firstPageResponse, err := client.GetThings(ctx, firstPageRequest)
	require.NoError(t, err)

	secondPageRequest := &boardgamegeek.GetThingsRequest{
		IDs:      []int64{catanId},
		Comments: true,
		Page:     2,
		PageSize: 10,
	}
	secondPageResponse, err := client.GetThings(ctx, secondPageRequest)
	require.NoError(t, err)

	assert.NotEqual(t,
		firstPageResponse.Items[0].Comments.Comments[0].Username,
		secondPageResponse.Items[0].Comments.Comments[0].Username,
	)
}

// Edge cases

func TestGetThingsByInvalidId(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs: []int64{999999999},
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.Empty(t, response.Items)
}

func TestGetThingsByMixOfValidAndInvalidIds(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs: []int64{catanId, 999999999},
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 1)
	assert.Equal(t, catanId, response.Items[0].Id)
}

func TestGetThingsOmitsOptionalSectionsByDefault(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs: []int64{catanId},
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)

	item := response.Items[0]
	assert.Nil(t, item.Versions)
	assert.Nil(t, item.Videos)
	assert.Nil(t, item.Statistics)
	assert.Nil(t, item.MarketplaceListings)
	assert.Nil(t, item.Comments)
}

func TestGetThingsPageAndPageSizeHaveNoEffectWithoutComments(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs:      []int64{catanId},
		Page:     1,
		PageSize: 5,
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.Nil(t, response.Items[0].Comments)
}

func TestGetThingsPageSizeBelowMinimumFallsBackToDefault(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.GetThingsRequest{
		IDs:      []int64{catanId},
		Comments: true,
		Page:     1,
		PageSize: 3,
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	// The BGG API's minimum pagesize is 10; anything smaller falls back to
	// the default page size of 100 rather than erroring or being clamped up to 10.
	assert.Len(t, response.Items[0].Comments.Comments, 100)
}
