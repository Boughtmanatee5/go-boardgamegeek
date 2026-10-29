//go:build integration

package integration_test

import (
	"testing"

	"github.com/Boughtmanatee5/go-boardgamegeek"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Catan (base game)
const catanID int64 = 13

// Catan: Seafarers (expansion)
const catanSeafarersID int64 = 325

func TestGetThingsByID(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ThingsOptions{
		IDs: []int64{catanID},
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 1)
	assert.Equal(t, catanID, response.Items[0].ID)
	assert.Equal(t, string(boardgamegeek.ThingTypeBoardGame), response.Items[0].Type)
	require.NotNil(t, response.Items[0].YearPublished)
	assert.Equal(t, "1995", string(*response.Items[0].YearPublished))
}

func TestGetThingsByMultipleIDs(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ThingsOptions{
		IDs: []int64{catanID, catanSeafarersID},
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 2)

	var returnedIDs []int64
	for _, item := range response.Items {
		returnedIDs = append(returnedIDs, item.ID)
	}
	assert.ElementsMatch(t, []int64{catanID, catanSeafarersID}, returnedIDs)
}

func TestGetThingsByType(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ThingsOptions{
		IDs:   []int64{catanID, catanSeafarersID},
		Types: []boardgamegeek.ThingType{boardgamegeek.ThingTypeBoardGameExpansion},
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 1)
	assert.Equal(t, catanSeafarersID, response.Items[0].ID)
	assert.Equal(t, string(boardgamegeek.ThingTypeBoardGameExpansion), response.Items[0].Type)
}

func TestGetThingsByRPGIssueType(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	// 140367 is an rpgissue: "& (Issue 1 - Summer 2012)"
	request := &boardgamegeek.ThingsOptions{
		IDs:   []int64{catanID, 140367},
		Types: []boardgamegeek.ThingType{boardgamegeek.ThingTypeRPGIssue},
	}

	response, err := client.GetThings(ctx, request)
	require.NoError(t, err)
	require.Len(t, response.Items, 1)
	assert.Equal(t, int64(140367), response.Items[0].ID)
	assert.Equal(t, string(boardgamegeek.ThingTypeRPGIssue), response.Items[0].Type)
}

func TestGetThingsWithVersions(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ThingsOptions{
		IDs:      []int64{catanID},
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
	request := &boardgamegeek.ThingsOptions{
		IDs:    []int64{catanID},
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
	request := &boardgamegeek.ThingsOptions{
		IDs:   []int64{catanID},
		Stats: true,
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response.Items[0].Statistics)
}

func TestGetThingsWithMarketPlace(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ThingsOptions{
		IDs:         []int64{catanID},
		MarketPlace: true,
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response.Items[0].MarketplaceListings)
}

func TestGetThingsWithComments(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ThingsOptions{
		IDs:      []int64{catanID},
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
	request := &boardgamegeek.ThingsOptions{
		IDs:            []int64{catanID},
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

	firstPageRequest := &boardgamegeek.ThingsOptions{
		IDs:      []int64{catanID},
		Comments: true,
		Page:     1,
		PageSize: 10,
	}
	firstPageResponse, err := client.GetThings(ctx, firstPageRequest)
	require.NoError(t, err)

	secondPageRequest := &boardgamegeek.ThingsOptions{
		IDs:      []int64{catanID},
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

func TestGetThingsByInvalidID(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ThingsOptions{
		IDs: []int64{999999999},
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.Empty(t, response.Items)
}

func TestGetThingsByMixOfValidAndInvalidIDs(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ThingsOptions{
		IDs: []int64{catanID, 999999999},
	}

	response, err := client.GetThings(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 1)
	assert.Equal(t, catanID, response.Items[0].ID)
}

func TestGetThingsOmitsOptionalSectionsByDefault(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ThingsOptions{
		IDs: []int64{catanID},
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
	request := &boardgamegeek.ThingsOptions{
		IDs:      []int64{catanID},
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
	request := &boardgamegeek.ThingsOptions{
		IDs:      []int64{catanID},
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
