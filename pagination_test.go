package boardgamegeek

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// pagesOf returns a fetchPage that serves items in pages of pageSize and records the pages requested
func pagesOf(items []int, pageSize int, total int, requested *[]int) fetchPage[int] {
	return func(_ context.Context, page int) ([]int, int, error) {
		*requested = append(*requested, page)
		start := min((page-1)*pageSize, len(items))
		end := min(start+pageSize, len(items))
		return items[start:end], total, nil
	}
}

func collect(seq func(yield func(int, error) bool)) ([]int, error) {
	var got []int
	for item, err := range seq {
		if err != nil {
			return got, err
		}
		got = append(got, item)
	}
	return got, nil
}

func TestPaginateStopsAtTotal(t *testing.T) {
	var requested []int
	got, err := collect(paginate(t.Context(), 0, pagesOf([]int{1, 2, 3, 4}, 2, 4, &requested)))

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, got)
	assert.Equal(t, []int{1, 2}, requested, "the total is reached on page 2, so page 3 is never requested")
}

func TestPaginateStopsAtShortPage(t *testing.T) {
	var requested []int
	got, err := collect(paginate(t.Context(), 0, pagesOf([]int{1, 2, 3}, 2, 0, &requested)))

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, got)
	assert.Equal(t, []int{1, 2}, requested)
}

func TestPaginateStopsAtEmptyPage(t *testing.T) {
	var requested []int
	got, err := collect(paginate(t.Context(), 0, pagesOf([]int{1, 2, 3, 4}, 2, 0, &requested)))

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, got)
	assert.Equal(t, []int{1, 2, 3}, requested, "without a total, a full last page needs one more request to find the end")
}

func TestPaginateFromStartPageCountsSkippedPagesTowardsTotal(t *testing.T) {
	var requested []int
	got, err := collect(paginate(t.Context(), 2, pagesOf([]int{1, 2, 3, 4, 5, 6}, 2, 6, &requested)))

	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4, 5, 6}, got)
	assert.Equal(t, []int{2, 3}, requested)
}

func TestPaginateYieldsErrorAndStops(t *testing.T) {
	fetchErr := errors.New("boom")
	calls := 0
	got, err := collect(paginate(t.Context(), 0, func(context.Context, int) ([]int, int, error) {
		calls++
		if calls == 2 {
			return nil, 0, fetchErr
		}
		return []int{calls, calls}, 10, nil
	}))

	assert.ErrorIs(t, err, fetchErr)
	assert.Equal(t, []int{1, 1}, got)
	assert.Equal(t, 2, calls)
}

func TestPaginateStopsWhenLoopBreaks(t *testing.T) {
	var requested []int
	for item, err := range paginate(t.Context(), 0, pagesOf([]int{1, 2, 3, 4}, 2, 4, &requested)) {
		require.NoError(t, err)
		if item == 1 {
			break
		}
	}
	assert.Equal(t, []int{1}, requested)
}

func TestAllThingCommentsNeedsOneID(t *testing.T) {
	c := NewClient()
	for _, err := range c.AllThingComments(t.Context(), &ThingsOptions{IDs: []int64{1, 2}}) {
		assert.ErrorIs(t, err, errThingCommentsID)
	}
}
