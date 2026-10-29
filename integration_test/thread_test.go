//go:build integration

package integration_test

import (
	"testing"

	"github.com/Boughtmanatee5/go-boardgamegeek"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// rulesForumThreads caches the first page of the Catan rules forum so each test
// doesn't spend a request on it; BGG rate limits the suite.
var rulesForumThreads []*boardgamegeek.ThreadSummary

// findThreadID returns a thread from the Catan rules forum with at least minArticles articles.
func findThreadID(t *testing.T, minArticles int) int64 {
	t.Helper()

	if rulesForumThreads == nil {
		client := newTestClient()
		forum, err := client.GetForum(t.Context(), &boardgamegeek.ForumOptions{ID: catanRulesForumID})
		require.NoError(t, err)
		rulesForumThreads = forum.Threads
	}

	for _, thread := range rulesForumThreads {
		if thread.NumArticles >= minArticles {
			return thread.ID
		}
	}
	t.Fatalf("no thread in forum %d with at least %d articles", catanRulesForumID, minArticles)
	return 0
}

func TestGetThreadByID(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	threadID := findThreadID(t, 1)

	response, err := client.GetThread(ctx, &boardgamegeek.ThreadOptions{ID: threadID})
	require.NoError(t, err)
	assert.Equal(t, threadID, response.ID)
	assert.NotEmpty(t, response.Articles)
}

func TestGetThreadPopulatesFields(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	threadID := findThreadID(t, 1)

	response, err := client.GetThread(ctx, &boardgamegeek.ThreadOptions{ID: threadID})
	require.NoError(t, err)
	assert.NotEmpty(t, response.Subject)
	assert.NotEmpty(t, response.Link)
	assert.NotEmpty(t, response.TermsOfUse)
	assert.Positive(t, response.NumArticles)

	article := response.Articles[0]
	assert.NotZero(t, article.ID)
	assert.NotEmpty(t, article.Username)
	assert.NotEmpty(t, article.Link)
	assert.NotEmpty(t, article.PostDate)
	assert.NotEmpty(t, article.EditDate)
	assert.NotEmpty(t, article.Subject)
	assert.NotEmpty(t, article.Body)
}

func TestGetThreadCountLimitsArticles(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	threadID := findThreadID(t, 3)

	response, err := client.GetThread(ctx, &boardgamegeek.ThreadOptions{ID: threadID, Count: 2})
	require.NoError(t, err)
	assert.Len(t, response.Articles, 2)
	assert.Greater(t, response.NumArticles, 2)
}

func TestGetThreadMinArticleID(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	threadID := findThreadID(t, 3)

	full, err := client.GetThread(ctx, &boardgamegeek.ThreadOptions{ID: threadID})
	require.NoError(t, err)
	minArticleID := full.Articles[1].ID

	response, err := client.GetThread(ctx, &boardgamegeek.ThreadOptions{ID: threadID, MinArticleID: minArticleID})
	require.NoError(t, err)
	require.NotEmpty(t, response.Articles)
	assert.Equal(t, minArticleID, response.Articles[0].ID)
	for _, article := range response.Articles {
		assert.GreaterOrEqual(t, article.ID, minArticleID)
	}
}

func TestGetThreadMinArticleDate(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	threadID := findThreadID(t, 1)

	past, err := client.GetThread(ctx, &boardgamegeek.ThreadOptions{ID: threadID, MinArticleDate: "2000-01-01"})
	require.NoError(t, err)
	assert.NotEmpty(t, past.Articles)

	future, err := client.GetThread(ctx, &boardgamegeek.ThreadOptions{ID: threadID, MinArticleDate: "2099-01-01 00:00:00"})
	require.NoError(t, err)
	assert.Equal(t, threadID, future.ID)
	assert.Empty(t, future.Articles)
}

func TestGetThreadUnknownIdReturnsError(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetThread(ctx, &boardgamegeek.ThreadOptions{ID: unknownID})
	assert.ErrorContains(t, err, "Thread Not Found")
	assert.ErrorIs(t, err, boardgamegeek.ErrNotFound)
	assert.Nil(t, response)
}
