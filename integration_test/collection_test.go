package integration_test

import (
	"errors"
	"testing"

	"github.com/Boughtmanatee5/go-boardgamegeek"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The 7th Continent, owned by Aldie
const seventhContinentID int64 = 180263

const collectionRetries = 10

// getCollection retries while BGG queues the request, as callers of GetCollection are expected to.
// The test client's transport spaces out the retries.
func getCollection(t *testing.T, request *boardgamegeek.CollectionOptions) (*boardgamegeek.CollectionResponse, error) {
	t.Helper()
	client := newTestClient()

	for range collectionRetries {
		response, err := client.GetCollection(t.Context(), request)
		if !errors.Is(err, boardgamegeek.ErrCollectionNotReady) {
			return response, err
		}
	}
	t.Fatalf("collection for %s still queued after %d attempts", request.Username, collectionRetries)
	return nil, nil
}

func TestGetCollectionByID(t *testing.T) {
	response, err := getCollection(t, &boardgamegeek.CollectionOptions{
		Username: aldieUsername,
		IDs:      []int64{seventhContinentID},
	})
	require.NoError(t, err)
	assert.Equal(t, 1, response.TotalItems)
	require.Len(t, response.Items, 1)
	assert.Equal(t, seventhContinentID, response.Items[0].ObjectID)
}

func TestGetCollectionPopulatesFields(t *testing.T) {
	response, err := getCollection(t, &boardgamegeek.CollectionOptions{
		Username: aldieUsername,
		IDs:      []int64{seventhContinentID},
	})
	require.NoError(t, err)
	assert.NotEmpty(t, response.TermsOfUse)
	assert.NotEmpty(t, response.PubDate)
	require.Len(t, response.Items, 1)

	item := response.Items[0]
	assert.Equal(t, "thing", item.ObjectType)
	assert.Equal(t, string(boardgamegeek.ThingTypeBoardGame), item.Subtype)
	assert.NotZero(t, item.CollID)
	require.NotNil(t, item.Name)
	assert.Equal(t, "The 7th Continent", item.Name.Value)
	assert.Equal(t, 5, item.Name.SortIndex)
	assert.Equal(t, "2017", item.YearPublished)
	assert.NotEmpty(t, item.Image)
	assert.NotEmpty(t, item.Thumbnail)
	assert.Positive(t, item.NumPlays)
	assert.NotEmpty(t, item.Comment)
	require.NotNil(t, item.Status)
	assert.True(t, item.Status.Own)
	assert.NotEmpty(t, item.Status.LastModified)
	assert.Nil(t, item.Stats)
	assert.Nil(t, item.Version)
}

func TestGetCollectionWithStats(t *testing.T) {
	response, err := getCollection(t, &boardgamegeek.CollectionOptions{
		Username: aldieUsername,
		IDs:      []int64{seventhContinentID},
		Stats:    true,
	})
	require.NoError(t, err)
	require.Len(t, response.Items, 1)

	stats := response.Items[0].Stats
	require.NotNil(t, stats)
	assert.Equal(t, 1, stats.MinPlayers)
	assert.Equal(t, 4, stats.MaxPlayers)
	assert.Positive(t, stats.PlayingTime)
	assert.Positive(t, stats.NumOwned)
	require.NotNil(t, stats.Rating)
	assert.NotEmpty(t, stats.Rating.Value)
	assert.NotEmpty(t, stats.Rating.Average)
	assert.NotEmpty(t, stats.Rating.UsersRated)
	require.NotEmpty(t, stats.Rating.Ranks)
	assert.Equal(t, "boardgame", stats.Rating.Ranks[0].Name)
}

func TestGetCollectionWithVersion(t *testing.T) {
	response, err := getCollection(t, &boardgamegeek.CollectionOptions{
		Username: aldieUsername,
		IDs:      []int64{seventhContinentID},
		Version:  true,
	})
	require.NoError(t, err)
	require.Len(t, response.Items, 1)
	require.NotNil(t, response.Items[0].Version)
	assert.Equal(t, "boardgameversion", response.Items[0].Version.Type)
	assert.NotZero(t, response.Items[0].Version.ID)
}

func TestGetCollectionOwnFilterIncludes(t *testing.T) {
	response, err := getCollection(t, &boardgamegeek.CollectionOptions{
		Username: aldieUsername,
		IDs:      []int64{seventhContinentID},
		Own:      boardgamegeek.Bool(true),
	})
	require.NoError(t, err)
	assert.Len(t, response.Items, 1)
}

func TestGetCollectionOwnFilterExcludes(t *testing.T) {
	response, err := getCollection(t, &boardgamegeek.CollectionOptions{
		Username: aldieUsername,
		IDs:      []int64{seventhContinentID},
		Own:      boardgamegeek.Bool(false),
	})
	require.NoError(t, err)
	assert.Empty(t, response.Items)
}

func TestGetCollectionBriefWishlist(t *testing.T) {
	response, err := getCollection(t, &boardgamegeek.CollectionOptions{
		Username: aldieUsername,
		Wishlist: boardgamegeek.Bool(true),
		Brief:    true,
	})
	require.NoError(t, err)
	require.NotEmpty(t, response.Items)
	for _, item := range response.Items {
		require.NotNil(t, item.Status)
		assert.True(t, item.Status.Wishlist)
		assert.Positive(t, item.Status.WishlistPriority)
		// brief omits everything but the name and status
		assert.Empty(t, item.Thumbnail)
	}
}

func TestGetCollectionUnknownUsernameReturnsError(t *testing.T) {
	response, err := getCollection(t, &boardgamegeek.CollectionOptions{
		Username: "zzznosuchuser123xyz",
	})
	assert.ErrorContains(t, err, "Invalid username specified")
	assert.ErrorIs(t, err, boardgamegeek.ErrNotFound)
	assert.Nil(t, response)
}
