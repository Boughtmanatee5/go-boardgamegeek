package integration_test

import (
	"testing"

	"github.com/Boughtmanatee5/go-boardgamegeek"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Catan: Rules forum
const catanRulesForumID int64 = 1508

// An id with no forum or thread behind it
const unknownID int64 = 999999999

func TestGetForumListByThing(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ForumListOptions{
		ID:   catanID,
		Type: boardgamegeek.ForumListTypeThing,
	}

	response, err := client.GetForumList(ctx, request)
	require.NoError(t, err)
	assert.Equal(t, catanID, response.ID)
	assert.Equal(t, string(boardgamegeek.ForumListTypeThing), response.Type)
	assert.NotEmpty(t, response.Forums)

	var forumIDs []int64
	for _, forum := range response.Forums {
		forumIDs = append(forumIDs, forum.ID)
	}
	assert.Contains(t, forumIDs, catanRulesForumID)
}

func TestGetForumListByFamily(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ForumListOptions{
		ID:   catanFamilyID,
		Type: boardgamegeek.ForumListTypeFamily,
	}

	response, err := client.GetForumList(ctx, request)
	require.NoError(t, err)
	assert.Equal(t, catanFamilyID, response.ID)
	assert.Equal(t, string(boardgamegeek.ForumListTypeFamily), response.Type)
	assert.NotEmpty(t, response.Forums)
}

func TestGetForumListPopulatesFields(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ForumListOptions{
		ID:   catanID,
		Type: boardgamegeek.ForumListTypeThing,
	}

	response, err := client.GetForumList(ctx, request)
	require.NoError(t, err)
	assert.NotEmpty(t, response.TermsOfUse)

	var rules *boardgamegeek.ForumSummary
	for _, forum := range response.Forums {
		if forum.ID == catanRulesForumID {
			rules = forum
		}
	}
	require.NotNil(t, rules)
	assert.Equal(t, "Rules", rules.Title)
	assert.NotEmpty(t, rules.Description)
	assert.False(t, rules.NoPosting)
	assert.Positive(t, rules.NumThreads)
	assert.Positive(t, rules.NumPosts)
	assert.NotEmpty(t, rules.LastPostDate)
}

func TestGetForumListUnknownIdReturnsEmptyForums(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ForumListOptions{
		ID:   unknownID,
		Type: boardgamegeek.ForumListTypeThing,
	}

	// BGG returns a default set of forums with no activity for unknown ids
	response, err := client.GetForumList(ctx, request)
	require.NoError(t, err)
	for _, forum := range response.Forums {
		assert.Zero(t, forum.NumThreads)
		assert.Zero(t, forum.NumPosts)
	}
}

func TestGetForumByID(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ForumOptions{
		ID: catanRulesForumID,
	}

	response, err := client.GetForum(ctx, request)
	require.NoError(t, err)
	assert.Equal(t, catanRulesForumID, response.ID)
	assert.Equal(t, "Rules", response.Title)
	assert.NotEmpty(t, response.Threads)
}

func TestGetForumPopulatesFields(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ForumOptions{
		ID: catanRulesForumID,
	}

	response, err := client.GetForum(ctx, request)
	require.NoError(t, err)
	assert.Positive(t, response.NumThreads)
	assert.Positive(t, response.NumPosts)
	assert.NotEmpty(t, response.TermsOfUse)

	thread := response.Threads[0]
	assert.NotZero(t, thread.ID)
	assert.NotEmpty(t, thread.Subject)
	assert.NotEmpty(t, thread.Author)
	assert.Positive(t, thread.NumArticles)
	assert.NotEmpty(t, thread.PostDate)
	assert.NotEmpty(t, thread.LastPostDate)
}

func TestGetForumPaging(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	pageOne, err := client.GetForum(ctx, &boardgamegeek.ForumOptions{ID: catanRulesForumID, Page: 1})
	require.NoError(t, err)
	pageTwo, err := client.GetForum(ctx, &boardgamegeek.ForumOptions{ID: catanRulesForumID, Page: 2})
	require.NoError(t, err)

	require.NotEmpty(t, pageOne.Threads)
	require.NotEmpty(t, pageTwo.Threads)
	assert.NotEqual(t, pageOne.Threads[0].ID, pageTwo.Threads[0].ID)
}

func TestGetForumPagePastEndIsEmpty(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ForumOptions{
		ID:   catanRulesForumID,
		Page: 999,
	}

	response, err := client.GetForum(ctx, request)
	require.NoError(t, err)
	assert.Equal(t, catanRulesForumID, response.ID)
	assert.Empty(t, response.Threads)
}

func TestGetForumUnknownIdReturnsError(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	request := &boardgamegeek.ForumOptions{
		ID: unknownID,
	}

	response, err := client.GetForum(ctx, request)
	assert.ErrorIs(t, err, boardgamegeek.ErrNotFound)
	assert.Nil(t, response)
}
