//go:build integration

package integration_test

import (
	"slices"
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
	assert.EqualValues(t, 1995, *response.Items[0].YearPublished)
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
	require.NotNil(t, response.Items[0].Versions)
	require.NotEmpty(t, response.Items[0].Versions.Items)

	for _, version := range response.Items[0].Versions.Items {
		assert.Equal(t, "boardgameversion", version.Type)
		assert.NotEmpty(t, version.Names)
		require.NotEmpty(t, version.Links)
		// every version links back to the thing it is a version of
		assert.True(t, slices.ContainsFunc(version.Links, func(link *boardgamegeek.Link) bool {
			return link.Type == "boardgameversion" && link.ID == catanID && link.Inbound
		}), "version %d links to Catan", version.ID)
	}
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
	require.NotNil(t, response.Items[0].Videos)
	assert.Greater(t, response.Items[0].Videos.Total, int64(0))
	require.NotEmpty(t, response.Items[0].Videos.Videos)

	for _, video := range response.Items[0].Videos.Videos {
		assert.Positive(t, video.ID)
		assert.NotEmpty(t, video.Title)
		assert.NotEmpty(t, video.Link)
		assert.False(t, video.PostDate.IsZero())
	}
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
	require.NotNil(t, response.Items[0].Statistics)
	ratings := response.Items[0].Statistics.Ratings
	require.NotNil(t, ratings)
	require.NotNil(t, ratings.UsersRated)
	assert.Positive(t, int(*ratings.UsersRated))
	require.NotNil(t, ratings.Average)
	assert.InDelta(t, 7, float64(*ratings.Average), 1.5)
	require.NotNil(t, ratings.AverageWeight)
	assert.InDelta(t, 2.5, float64(*ratings.AverageWeight), 1.5)

	require.NotEmpty(t, ratings.Ranks)
	assert.Equal(t, "subtype", ratings.Ranks[0].Type)
	assert.Equal(t, "boardgame", ratings.Ranks[0].Name)
	assert.NotEmpty(t, ratings.Ranks[0].Value)
}

func TestGetThingsPopulatesBoardGameDetails(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetThings(ctx, &boardgamegeek.ThingsOptions{IDs: []int64{catanID}})
	require.NoError(t, err)
	require.Len(t, response.Items, 1)
	item := response.Items[0]

	assert.NotEmpty(t, item.Thumbnail)
	assert.NotEmpty(t, item.Image)
	assert.NotEmpty(t, item.Description)
	assert.Equal(t, boardgamegeek.IntValue(3), *item.MinPlayers)
	assert.Equal(t, boardgamegeek.IntValue(4), *item.MaxPlayers)
	require.NotNil(t, item.PlayingTime)
	assert.Positive(t, int(*item.PlayingTime))
	require.NotNil(t, item.MinAge)
	assert.Positive(t, int(*item.MinAge))

	linkTypes := map[string]bool{}
	for _, link := range item.Links {
		linkTypes[link.Type] = true
	}
	for _, linkType := range []string{"boardgamecategory", "boardgamemechanic", "boardgamedesigner", "boardgamepublisher", "boardgameexpansion"} {
		assert.True(t, linkTypes[linkType], "Catan has a %s link", linkType)
	}

	pollNames := map[string]bool{}
	for _, poll := range item.Polls {
		pollNames[poll.Name] = true
		assert.Positive(t, poll.TotalVotes)
		assert.NotEmpty(t, poll.Results)
	}
	assert.Equal(t, map[string]bool{"suggested_numplayers": true, "suggested_playerage": true, "language_dependence": true}, pollNames)
	require.NotNil(t, item.PollSummary)
	assert.NotEmpty(t, item.PollSummary.Results)
}

func TestGetThingsPopulatesRPGIssueDetails(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	// 140367 is an rpgissue: "& (Issue 1 - Summer 2012)"
	response, err := client.GetThings(ctx, &boardgamegeek.ThingsOptions{IDs: []int64{140367}, Versions: true})
	require.NoError(t, err)
	require.Len(t, response.Items, 1)
	item := response.Items[0]

	require.NotNil(t, item.DatePublished)
	assert.Equal(t, "2012-07-01", string(*item.DatePublished))
	assert.Equal(t, boardgamegeek.IntValue(1), *item.IssueIndex)
	assert.Nil(t, item.MinPlayers)
	assert.Empty(t, item.Polls)

	require.NotNil(t, item.Versions)
	require.NotEmpty(t, item.Versions.Items)
	version := item.Versions.Items[0]
	require.NotNil(t, version.PageCount)
	assert.Positive(t, int(*version.PageCount))
	assert.NotNil(t, version.Format)
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
	require.NotNil(t, response.Items[0].MarketplaceListings)
	require.NotEmpty(t, response.Items[0].MarketplaceListings.Listings)

	for _, listing := range response.Items[0].MarketplaceListings.Listings {
		require.NotNil(t, listing.Price)
		assert.NotEmpty(t, listing.Price.Currency)
		assert.NotEmpty(t, listing.Price.Value)
		require.NotNil(t, listing.Link)
		assert.NotEmpty(t, listing.Link.Href)
	}
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
