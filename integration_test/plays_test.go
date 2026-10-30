//go:build integration

package integration_test

import (
	"slices"
	"testing"
	"time"

	"github.com/Boughtmanatee5/go-boardgamegeek"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BGG returns plays 100 to a page
const playsPageSize = 100

func TestGetPlaysByUser(t *testing.T) {
	client := newTestClient()

	response, err := client.GetPlays(t.Context(), &boardgamegeek.PlaysOptions{Username: aldieUsername})
	require.NoError(t, err)
	assert.Equal(t, aldieUsername, response.Username)
	assert.Equal(t, aldieID, response.UserID)
	assert.Equal(t, 1, response.Page)
	assert.Greater(t, response.Total, playsPageSize)
	require.Len(t, response.Plays, playsPageSize)

	for _, play := range response.Plays {
		assert.Positive(t, play.ID)
		require.NotNil(t, play.Item)
		assert.NotEmpty(t, play.Item.Name)
	}
}

func TestGetPlaysByThing(t *testing.T) {
	client := newTestClient()

	response, err := client.GetPlays(t.Context(), &boardgamegeek.PlaysOptions{ID: catanID, Type: boardgamegeek.PlaysTypeThing})
	require.NoError(t, err)
	assert.Empty(t, response.Username)
	require.NotEmpty(t, response.Plays)

	for _, play := range response.Plays {
		// plays listed by thing carry the user who logged each one
		assert.Positive(t, play.UserID)
		require.NotNil(t, play.Item)
		assert.Equal(t, catanID, play.Item.ObjectID)
	}
}

func TestGetPlaysByUserAndThing(t *testing.T) {
	client := newTestClient()

	response, err := client.GetPlays(t.Context(), &boardgamegeek.PlaysOptions{
		Username: aldieUsername,
		ID:       seventhContinentID,
		Type:     boardgamegeek.PlaysTypeThing,
	})
	require.NoError(t, err)
	require.NotEmpty(t, response.Plays)
	assert.Len(t, response.Plays, response.Total)

	for _, play := range response.Plays {
		require.NotNil(t, play.Item)
		assert.Equal(t, seventhContinentID, play.Item.ObjectID)
	}
}

func TestGetPlaysByDateRange(t *testing.T) {
	client := newTestClient()

	response, err := client.GetPlays(t.Context(), &boardgamegeek.PlaysOptions{
		Username: aldieUsername,
		MinDate:  time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
		MaxDate:  time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	require.NotEmpty(t, response.Plays)
	assert.Less(t, response.Total, playsPageSize)

	for _, play := range response.Plays {
		assert.GreaterOrEqual(t, play.Date, "2026-08-01")
		assert.LessOrEqual(t, play.Date, "2026-10-31")
	}
}

func TestGetPlaysBySubtype(t *testing.T) {
	client := newTestClient()

	response, err := client.GetPlays(t.Context(), &boardgamegeek.PlaysOptions{
		Username: aldieUsername,
		Subtype:  boardgamegeek.ThingTypeBoardGameExpansion,
	})
	require.NoError(t, err)
	require.NotEmpty(t, response.Plays)

	for _, play := range response.Plays {
		require.NotNil(t, play.Item)
		assert.True(t, slices.Contains(play.Item.Subtypes, boardgamegeek.ValueAttr(boardgamegeek.ThingTypeBoardGameExpansion)),
			"play %d of %s is an expansion", play.ID, play.Item.Name)
	}
}

// Edge cases

func TestGetPlaysPastLastPageIsEmpty(t *testing.T) {
	client := newTestClient()

	response, err := client.GetPlays(t.Context(), &boardgamegeek.PlaysOptions{Username: aldieUsername, Page: 9999})
	require.NoError(t, err)
	assert.Positive(t, response.Total)
	assert.Empty(t, response.Plays)
}

func TestGetPlaysUnknownUserIsNotFound(t *testing.T) {
	client := newTestClient()

	_, err := client.GetPlays(t.Context(), &boardgamegeek.PlaysOptions{Username: "nosuchuser_zzq9"})
	assert.ErrorIs(t, err, boardgamegeek.ErrNotFound)
}

func TestGetPlaysWithoutUserOrIDIsAnError(t *testing.T) {
	client := newTestClient()

	_, err := client.GetPlays(t.Context(), &boardgamegeek.PlaysOptions{})
	require.Error(t, err)
	assert.NotErrorIs(t, err, boardgamegeek.ErrNotFound)
}
