package integration_test

import (
	"testing"

	boardgamegeek "github.com/Boughtmanatee5/BoardGameGeekClient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// rulesForumThreads caches the first page of the Catan rules forum so each test
// doesn't spend a request on it; BGG rate limits the suite.
var rulesForumThreads []*boardgamegeek.ThreadSummary

// findThreadId returns a thread from the Catan rules forum with at least minArticles articles.
func findThreadId(t *testing.T, minArticles int) int64 {
	t.Helper()

	if rulesForumThreads == nil {
		client := newTestClient()
		forum, err := client.GetForum(t.Context(), &boardgamegeek.ForumRequest{ID: catanRulesForumId})
		require.NoError(t, err)
		rulesForumThreads = forum.Threads
	}

	for _, thread := range rulesForumThreads {
		if thread.NumArticles >= minArticles {
			return thread.Id
		}
	}
	t.Fatalf("no thread in forum %d with at least %d articles", catanRulesForumId, minArticles)
	return 0
}

func TestGetThreadById(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	threadId := findThreadId(t, 1)

	response, err := client.GetThread(ctx, &boardgamegeek.ThreadRequest{ID: threadId})
	require.NoError(t, err)
	assert.Equal(t, threadId, response.Id)
	assert.NotEmpty(t, response.Articles)
}

func TestGetThreadPopulatesFields(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	threadId := findThreadId(t, 1)

	response, err := client.GetThread(ctx, &boardgamegeek.ThreadRequest{ID: threadId})
	require.NoError(t, err)
	assert.NotEmpty(t, response.Subject)
	assert.NotEmpty(t, response.Link)
	assert.NotEmpty(t, response.TermsOfUse)
	assert.Positive(t, response.NumArticles)

	article := response.Articles[0]
	assert.NotZero(t, article.Id)
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
	threadId := findThreadId(t, 3)

	response, err := client.GetThread(ctx, &boardgamegeek.ThreadRequest{ID: threadId, Count: 2})
	require.NoError(t, err)
	assert.Len(t, response.Articles, 2)
	assert.Greater(t, response.NumArticles, 2)
}

func TestGetThreadMinArticleId(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	threadId := findThreadId(t, 3)

	full, err := client.GetThread(ctx, &boardgamegeek.ThreadRequest{ID: threadId})
	require.NoError(t, err)
	minArticleId := full.Articles[1].Id

	response, err := client.GetThread(ctx, &boardgamegeek.ThreadRequest{ID: threadId, MinArticleID: minArticleId})
	require.NoError(t, err)
	require.NotEmpty(t, response.Articles)
	assert.Equal(t, minArticleId, response.Articles[0].Id)
	for _, article := range response.Articles {
		assert.GreaterOrEqual(t, article.Id, minArticleId)
	}
}

func TestGetThreadMinArticleDate(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()
	threadId := findThreadId(t, 1)

	past, err := client.GetThread(ctx, &boardgamegeek.ThreadRequest{ID: threadId, MinArticleDate: "2000-01-01"})
	require.NoError(t, err)
	assert.NotEmpty(t, past.Articles)

	future, err := client.GetThread(ctx, &boardgamegeek.ThreadRequest{ID: threadId, MinArticleDate: "2099-01-01 00:00:00"})
	require.NoError(t, err)
	assert.Equal(t, threadId, future.Id)
	assert.Empty(t, future.Articles)
}

func TestGetThreadUnknownIdReturnsError(t *testing.T) {
	ctx := t.Context()
	client := newTestClient()

	response, err := client.GetThread(ctx, &boardgamegeek.ThreadRequest{ID: unknownId})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "Thread Not Found")
	assert.Nil(t, response)
}
