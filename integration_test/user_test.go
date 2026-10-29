package integration_test

import (
	"net/http"
	"testing"

	"github.com/Boughtmanatee5/go-boardgamegeek"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Scott Alden, BGG co-founder
const aldieUsername = "Aldie"
const aldieID int64 = 688

func TestGetUserByName(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetUser(ctx, &boardgamegeek.UserOptions{Name: aldieUsername})
	require.NoError(t, err)
	assert.Equal(t, aldieID, response.ID)
	assert.Equal(t, aldieUsername, response.Name)
}

func TestGetUserPopulatesProfile(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetUser(ctx, &boardgamegeek.UserOptions{Name: aldieUsername})
	require.NoError(t, err)
	assert.NotEmpty(t, response.TermsOfUse)
	assert.Equal(t, "Scott", string(response.FirstName))
	assert.Equal(t, "Alden", string(response.LastName))
	assert.Equal(t, "1999", string(response.YearRegistered))
	assert.Equal(t, "United States", string(response.Country))
	assert.NotEmpty(t, response.AvatarLink)
	assert.NotEmpty(t, response.LastLogin)
	assert.NotEmpty(t, response.TradeRating)
}

func TestGetUserOmitsOptionalSectionsByDefault(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetUser(ctx, &boardgamegeek.UserOptions{Name: aldieUsername})
	require.NoError(t, err)
	assert.Nil(t, response.Buddies)
	assert.Nil(t, response.Guilds)
	assert.Nil(t, response.Hot)
	assert.Nil(t, response.Top)
}

func TestGetUserWithBuddies(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetUser(ctx, &boardgamegeek.UserOptions{Name: aldieUsername, Buddies: true})
	require.NoError(t, err)
	require.NotNil(t, response.Buddies)
	assert.Positive(t, response.Buddies.Total)
	assert.Equal(t, 1, response.Buddies.Page)
	require.NotEmpty(t, response.Buddies.Buddies)
	assert.NotZero(t, response.Buddies.Buddies[0].ID)
	assert.NotEmpty(t, response.Buddies.Buddies[0].Name)
}

func TestGetUserWithGuilds(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetUser(ctx, &boardgamegeek.UserOptions{Name: aldieUsername, Guilds: true})
	require.NoError(t, err)
	require.NotNil(t, response.Guilds)
	assert.Positive(t, response.Guilds.Total)

	var guildIDs []int64
	for _, guild := range response.Guilds.Guilds {
		guildIDs = append(guildIDs, guild.ID)
	}
	assert.Contains(t, guildIDs, dallasMetrogamersGuildID)
}

func TestGetUserWithHotAndTop(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.UserOptions{
		Name: aldieUsername,
		Hot:  true,
		Top:  true,
	}

	response, err := client.GetUser(ctx, request)
	require.NoError(t, err)
	require.NotNil(t, response.Hot)
	require.NotNil(t, response.Top)
	assert.Equal(t, string(boardgamegeek.UserDomainBoardGame), response.Hot.Domain)
	assert.Equal(t, string(boardgamegeek.UserDomainBoardGame), response.Top.Domain)
	require.NotEmpty(t, response.Top.Items)

	item := response.Top.Items[0]
	assert.Equal(t, 1, item.Rank)
	assert.Equal(t, "thing", item.Type)
	assert.NotZero(t, item.ID)
	assert.NotEmpty(t, item.Name)
}

func TestGetUserTopByDomain(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.UserOptions{
		Name:   aldieUsername,
		Top:    true,
		Domain: boardgamegeek.UserDomainRPG,
	}

	response, err := client.GetUser(ctx, request)
	require.NoError(t, err)
	require.NotNil(t, response.Top)
	assert.Equal(t, string(boardgamegeek.UserDomainRPG), response.Top.Domain)
}

func TestGetUserUnknownNameReturnsError(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetUser(ctx, &boardgamegeek.UserOptions{Name: "zzznosuchuser123xyz"})
	assert.ErrorIs(t, err, boardgamegeek.ErrNotFound)
	var errResp *boardgamegeek.ErrorResponse
	require.ErrorAs(t, err, &errResp)
	assert.Equal(t, http.StatusNotFound, errResp.Response.StatusCode)
	assert.Nil(t, response)
}
