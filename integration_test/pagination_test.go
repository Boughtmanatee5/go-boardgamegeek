//go:build integration

package integration_test

import (
	"testing"

	"github.com/Boughtmanatee5/go-boardgamegeek"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// take collects up to n items from seq, failing the test on the first error
func take[T any](t *testing.T, seq func(yield func(T, error) bool), n int) []T {
	t.Helper()
	var items []T
	for item, err := range seq {
		require.NoError(t, err)
		items = append(items, item)
		if len(items) == n {
			break
		}
	}
	return items
}

func TestAllForumThreadsSpansPages(t *testing.T) {
	client := newTestClient()

	threads := take(t, client.AllForumThreads(t.Context(), &boardgamegeek.ForumOptions{ID: catanRulesForumID}), 60)
	require.Len(t, threads, 60, "the forum has more than one page of threads")

	ids := map[int64]bool{}
	for _, thread := range threads {
		ids[thread.ID] = true
	}
	assert.Len(t, ids, 60, "no thread is repeated across pages")
}

func TestAllGuildMembersYieldsEveryMember(t *testing.T) {
	client := newTestClient()
	ctx := t.Context()

	firstPage, err := client.GetGuild(ctx, &boardgamegeek.GuildOptions{ID: dallasMetrogamersGuildID, Members: true})
	require.NoError(t, err)

	members := take(t, client.AllGuildMembers(ctx, &boardgamegeek.GuildOptions{ID: dallasMetrogamersGuildID}), -1)
	assert.Len(t, members, firstPage.Members.Count)
}

func TestAllUserGuildsYieldsEveryGuild(t *testing.T) {
	client := newTestClient()
	ctx := t.Context()

	firstPage, err := client.GetUser(ctx, &boardgamegeek.UserOptions{Name: aldieUsername, Guilds: true})
	require.NoError(t, err)

	guilds := take(t, client.AllUserGuilds(ctx, &boardgamegeek.UserOptions{Name: aldieUsername}), -1)
	assert.Len(t, guilds, firstPage.Guilds.Total)
}

func TestAllUserBuddiesYieldsEveryBuddy(t *testing.T) {
	client := newTestClient()
	ctx := t.Context()

	firstPage, err := client.GetUser(ctx, &boardgamegeek.UserOptions{Name: aldieUsername, Buddies: true})
	require.NoError(t, err)

	buddies := take(t, client.AllUserBuddies(ctx, &boardgamegeek.UserOptions{Name: aldieUsername}), -1)
	assert.Len(t, buddies, firstPage.Buddies.Total)

	ids := map[int64]bool{}
	for _, buddy := range buddies {
		ids[buddy.ID] = true
	}
	assert.Len(t, ids, firstPage.Buddies.Total, "no buddy is repeated")
}

func TestAllThingCommentsSpansPages(t *testing.T) {
	client := newTestClient()

	comments := take(t, client.AllThingComments(t.Context(), &boardgamegeek.ThingsOptions{IDs: []int64{catanID}, PageSize: 10}), 25)
	require.Len(t, comments, 25)

	usernames := map[string]bool{}
	for _, comment := range comments {
		usernames[comment.Username] = true
	}
	assert.Len(t, usernames, 25, "no comment is repeated across pages")
}
