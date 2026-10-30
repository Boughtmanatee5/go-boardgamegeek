package boardgamegeek

import (
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRequestQueryEncoding checks the query string each request struct's url tags produce.
// url.Values.Encode sorts parameters by key.
func TestRequestQueryEncoding(t *testing.T) {
	tests := []struct {
		name    string
		request any
		want    string
	}{
		{
			name:    "search escapes reserved characters",
			request: &SearchOptions{Query: "Ticket & Ride: Europe"},
			want:    "query=Ticket+%26+Ride%3A+Europe",
		},
		{
			name:    "search with types and exact",
			request: &SearchOptions{Query: "catan", Types: []SearchType{SearchTypeBoardGame, SearchTypeBoardGameExpansion}, Exact: true},
			want:    "exact=1&query=catan&type=boardgame%2Cboardgameexpansion",
		},
		{
			name:    "search with non-thing types",
			request: &SearchOptions{Query: "dragon", Types: []SearchType{SearchTypeBoardGameDesigner, SearchTypeRPGIssue}},
			want:    "query=dragon&type=boardgamedesigner%2Crpgissue",
		},
		{
			name:    "search keeps an empty query",
			request: &SearchOptions{},
			want:    "query=",
		},
		{
			name:    "things with flags and paging",
			request: &ThingsOptions{IDs: []int64{1, 2}, Stats: true, Comments: true, Page: 2, PageSize: 50},
			want:    "comments=1&id=1%2C2&page=2&pagesize=50&stats=1",
		},
		{
			name:    "things omits false flags and zero paging",
			request: &ThingsOptions{IDs: []int64{13}},
			want:    "id=13",
		},
		{
			name:    "family with types",
			request: &FamilyItemsOptions{IDs: []int64{3}, Types: []FamilyItemType{FamilyItemTypeBoardGameFamily}},
			want:    "id=3&type=boardgamefamily",
		},
		{
			name:    "forum list",
			request: &ForumListOptions{ID: 13, Type: ForumListTypeThing},
			want:    "id=13&type=thing",
		},
		{
			name:    "forum page",
			request: &ForumOptions{ID: 13, Page: 3},
			want:    "id=13&page=3",
		},
		{
			name:    "thread escapes the date",
			request: &ThreadOptions{ID: 5, MinArticleDate: "2020-01-01 10:00:00", Count: 2},
			want:    "count=2&id=5&minarticledate=2020-01-01+10%3A00%3A00",
		},
		{
			name:    "hot without type is empty",
			request: &HotItemsOptions{},
			want:    "",
		},
		{
			name:    "user with options",
			request: &UserOptions{Name: "some user", Buddies: true, Domain: UserDomainRPG, Page: 2},
			want:    "buddies=1&domain=rpg&name=some+user&page=2",
		},
		{
			name:    "guild with members",
			request: &GuildOptions{ID: 1, Members: true, Sort: GuildMemberSortDate},
			want:    "id=1&members=1&sort=date",
		},
		{
			name: "collection status filters encode nil, true and false",
			request: &CollectionOptions{
				Username:  "Aldie",
				Own:       new(true),
				Wishlist:  new(false),
				MinRating: 7.5,
				IDs:       []int64{1, 2},
			},
			want: "id=1%2C2&minrating=7.5&own=1&username=Aldie&wishlist=0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := query.Values(tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.want, params.Encode())
		})
	}
}
