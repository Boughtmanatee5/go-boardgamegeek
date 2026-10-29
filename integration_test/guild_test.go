package integration_test

import (
	"testing"

	boardgamegeek "github.com/Boughtmanatee5/BoardGameGeekClient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Dallas Metrogamers, managed by Aldie
const dallasMetrogamersGuildId int64 = 1

// BGG returns guild members 25 per page
const guildMembersPageSize = 25

func TestGetGuildById(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetGuild(ctx, &boardgamegeek.GuildRequest{ID: dallasMetrogamersGuildId})
	require.NoError(t, err)
	assert.Equal(t, dallasMetrogamersGuildId, response.Id)
	assert.Equal(t, "Dallas Metrogamers", response.Name)
}

func TestGetGuildPopulatesFields(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetGuild(ctx, &boardgamegeek.GuildRequest{ID: dallasMetrogamersGuildId})
	require.NoError(t, err)
	assert.NotEmpty(t, response.Created)
	assert.NotEmpty(t, response.TermsOfUse)
	assert.Equal(t, "group", response.Category)
	assert.Equal(t, aldieUsername, response.Manager)
	assert.NotEmpty(t, response.Description)
	require.NotNil(t, response.Location)
	assert.Equal(t, "Dallas", response.Location.City)
	assert.Equal(t, "United States", response.Location.Country)
}

func TestGetGuildOmitsMembersByDefault(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetGuild(ctx, &boardgamegeek.GuildRequest{ID: dallasMetrogamersGuildId})
	require.NoError(t, err)
	assert.Nil(t, response.Members)
}

func TestGetGuildWithMembers(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetGuild(ctx, &boardgamegeek.GuildRequest{ID: dallasMetrogamersGuildId, Members: true})
	require.NoError(t, err)
	require.NotNil(t, response.Members)
	assert.Greater(t, response.Members.Count, guildMembersPageSize)
	assert.Equal(t, 1, response.Members.Page)
	assert.Len(t, response.Members.Members, guildMembersPageSize)
	assert.NotEmpty(t, response.Members.Members[0].Name)
	assert.NotEmpty(t, response.Members.Members[0].Date)
}

func TestGetGuildMembersPaging(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	pageOne, err := client.GetGuild(ctx, &boardgamegeek.GuildRequest{ID: dallasMetrogamersGuildId, Members: true, Page: 1})
	require.NoError(t, err)
	pageTwo, err := client.GetGuild(ctx, &boardgamegeek.GuildRequest{ID: dallasMetrogamersGuildId, Members: true, Page: 2})
	require.NoError(t, err)

	assert.Equal(t, 2, pageTwo.Members.Page)
	require.NotEmpty(t, pageTwo.Members.Members)
	assert.NotEqual(t, pageOne.Members.Members[0].Name, pageTwo.Members.Members[0].Name)
}

func TestGetGuildMembersSort(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	byUsername, err := client.GetGuild(ctx, &boardgamegeek.GuildRequest{
		ID: dallasMetrogamersGuildId, Members: true, Sort: boardgamegeek.GuildMemberSortUsername,
	})
	require.NoError(t, err)
	byDate, err := client.GetGuild(ctx, &boardgamegeek.GuildRequest{
		ID: dallasMetrogamersGuildId, Members: true, Sort: boardgamegeek.GuildMemberSortDate,
	})
	require.NoError(t, err)

	var usernameOrder, dateOrder []string
	for _, member := range byUsername.Members.Members {
		usernameOrder = append(usernameOrder, member.Name)
	}
	for _, member := range byDate.Members.Members {
		dateOrder = append(dateOrder, member.Name)
	}
	assert.NotEqual(t, usernameOrder, dateOrder)
}

func TestGetGuildUnknownIdReturnsError(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetGuild(ctx, &boardgamegeek.GuildRequest{ID: unknownId})
	assert.ErrorContains(t, err, "Not Found")
	assert.ErrorIs(t, err, boardgamegeek.ErrNotFound)
	assert.Nil(t, response)
}
