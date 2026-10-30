package boardgamegeek_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Boughtmanatee5/go-boardgamegeek"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests parse real BGG responses saved in testdata/<endpoint>.xml, so they run offline.
// The request each fixture was captured with is noted on its test.

// newFixtureClient returns a Client whose requests to /<endpoint> are answered with testdata/<endpoint>.xml
func newFixtureClient(t *testing.T) *boardgamegeek.Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := os.ReadFile(filepath.Join("testdata", strings.TrimPrefix(r.URL.Path, "/")+".xml"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)

	baseURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	return boardgamegeek.NewClient(boardgamegeek.WithBaseURL(baseURL))
}

// search?query=catan&type=boardgame&exact=1
func TestSearchParsesFixture(t *testing.T) {
	resp, err := newFixtureClient(t).Search(t.Context(), &boardgamegeek.SearchOptions{Query: "catan"})
	require.NoError(t, err)

	assert.EqualValues(t, 1, resp.Total)
	require.Len(t, resp.Items, 1)
	assert.EqualValues(t, 13, resp.Items[0].ID)
	assert.Equal(t, "Catan", resp.Items[0].Name.Value)
	require.NotNil(t, resp.Items[0].YearPublished)
	assert.Equal(t, "1995", string(*resp.Items[0].YearPublished))
}

// thing?id=13&stats=1
func TestGetThingsParsesFixture(t *testing.T) {
	resp, err := newFixtureClient(t).GetThings(t.Context(), &boardgamegeek.ThingsOptions{IDs: []int64{13}})
	require.NoError(t, err)

	require.Len(t, resp.Items, 1)
	thing := resp.Items[0]
	assert.EqualValues(t, 13, thing.ID)
	assert.Equal(t, string(boardgamegeek.ThingTypeBoardGame), thing.Type)
	assert.NotEmpty(t, thing.Names)
	require.NotNil(t, thing.YearPublished)
	assert.Equal(t, "1995", string(*thing.YearPublished))
	assert.NotNil(t, thing.Statistics)
	assert.Nil(t, thing.Comments)
}

// family?id=3
func TestGetFamilyItemsParsesFixture(t *testing.T) {
	resp, err := newFixtureClient(t).GetFamilyItems(t.Context(), &boardgamegeek.FamilyItemsOptions{IDs: []int64{3}})
	require.NoError(t, err)

	require.Len(t, resp.Items, 1)
	assert.EqualValues(t, 3, resp.Items[0].ID)
	assert.Equal(t, string(boardgamegeek.FamilyItemTypeBoardGameFamily), resp.Items[0].Type)
	assert.Len(t, resp.Items[0].Links, 223)
}

// forumlist?id=13&type=thing
func TestGetForumListParsesFixture(t *testing.T) {
	resp, err := newFixtureClient(t).GetForumList(t.Context(), &boardgamegeek.ForumListOptions{ID: 13, Type: boardgamegeek.ForumListTypeThing})
	require.NoError(t, err)

	assert.EqualValues(t, 13, resp.ID)
	require.NotEmpty(t, resp.Forums)
	assert.EqualValues(t, 297, resp.Forums[0].ID)
	assert.Equal(t, "Reviews", resp.Forums[0].Title)
	assert.Positive(t, resp.Forums[0].NumThreads)
}

// forum?id=1508
func TestGetForumParsesFixture(t *testing.T) {
	resp, err := newFixtureClient(t).GetForum(t.Context(), &boardgamegeek.ForumOptions{ID: 1508})
	require.NoError(t, err)

	assert.EqualValues(t, 1508, resp.ID)
	assert.Len(t, resp.Threads, 50)
	assert.Positive(t, resp.Threads[0].ID)
}

// thread?id=3782789&count=2
func TestGetThreadParsesFixture(t *testing.T) {
	resp, err := newFixtureClient(t).GetThread(t.Context(), &boardgamegeek.ThreadOptions{ID: 3782789, Count: 2})
	require.NoError(t, err)

	assert.EqualValues(t, 3782789, resp.ID)
	assert.Equal(t, "Old rulebook", resp.Subject)
	require.Len(t, resp.Articles, 2)
	assert.EqualValues(t, 48294493, resp.Articles[0].ID)
	assert.NotEmpty(t, resp.Articles[0].Body)
}

// hot?type=boardgame
func TestGetHotItemsParsesFixture(t *testing.T) {
	resp, err := newFixtureClient(t).GetHotItems(t.Context(), &boardgamegeek.HotItemsOptions{})
	require.NoError(t, err)

	require.Len(t, resp.Items, 50)
	assert.Equal(t, 1, resp.Items[0].Rank)
	assert.NotEmpty(t, resp.Items[0].Name)
}

// user?name=Aldie&guilds=1
func TestGetUserParsesFixture(t *testing.T) {
	resp, err := newFixtureClient(t).GetUser(t.Context(), &boardgamegeek.UserOptions{Name: "Aldie", Guilds: true})
	require.NoError(t, err)

	assert.EqualValues(t, 688, resp.ID)
	assert.Equal(t, "Aldie", resp.Name)
	assert.Equal(t, "Scott", string(resp.FirstName))
	assert.Equal(t, "United States", string(resp.Country))
	require.NotNil(t, resp.Guilds)
	assert.Len(t, resp.Guilds.Guilds, resp.Guilds.Total)
	assert.Nil(t, resp.Buddies)
}

// guild?id=1&members=1
func TestGetGuildParsesFixture(t *testing.T) {
	resp, err := newFixtureClient(t).GetGuild(t.Context(), &boardgamegeek.GuildOptions{ID: 1, Members: true})
	require.NoError(t, err)

	assert.EqualValues(t, 1, resp.ID)
	assert.Equal(t, "Dallas Metrogamers", resp.Name)
	assert.Equal(t, "Aldie", resp.Manager)
	require.NotNil(t, resp.Location)
	require.NotNil(t, resp.Members)
	assert.Equal(t, 1, resp.Members.Page)
	assert.NotEmpty(t, resp.Members.Members)
}

// collection?username=Aldie&id=180263&stats=1
func TestGetCollectionParsesFixture(t *testing.T) {
	resp, err := newFixtureClient(t).GetCollection(t.Context(), &boardgamegeek.CollectionOptions{Username: "Aldie", IDs: []int64{180263}})
	require.NoError(t, err)

	assert.Equal(t, 1, resp.TotalItems)
	require.Len(t, resp.Items, 1)
	item := resp.Items[0]
	assert.EqualValues(t, 180263, item.ObjectID)
	require.NotNil(t, item.Name)
	assert.Equal(t, "The 7th Continent", item.Name.Value)
	assert.Equal(t, "2017", item.YearPublished)
	assert.NotNil(t, item.Stats)
	assert.NotNil(t, item.Status)
}

// plays?username=Aldie&mindate=2026-08-01&maxdate=2026-10-31
func TestGetPlaysParsesFixture(t *testing.T) {
	resp, err := newFixtureClient(t).GetPlays(t.Context(), &boardgamegeek.PlaysOptions{Username: "Aldie"})
	require.NoError(t, err)

	assert.Equal(t, "Aldie", resp.Username)
	assert.EqualValues(t, 688, resp.UserID)
	assert.Equal(t, 7, resp.Total)
	assert.Equal(t, 1, resp.Page)
	require.Len(t, resp.Plays, 7)

	play := resp.Plays[0]
	assert.EqualValues(t, 119593167, play.ID)
	assert.Equal(t, "2026-10-09", play.Date)
	assert.Equal(t, 5, play.Quantity)
	assert.False(t, play.Incomplete)
	require.NotNil(t, play.Item)
	assert.Equal(t, "Dice Conquest", play.Item.Name)
	assert.Equal(t, "thing", play.Item.ObjectType)
	assert.EqualValues(t, 354866, play.Item.ObjectID)
	assert.Equal(t, []boardgamegeek.ValueAttr{"boardgame"}, play.Item.Subtypes)
	assert.Empty(t, play.Players)
}

// testdata/plays_players.xml is a hand-written mock in the shape of plays?id=13&type=thing,
// since real responses by thing list other users' plays
func TestGetPlaysParsesPlayersMock(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("testdata", "plays_players.xml"))
	require.NoError(t, err)
	client := newStubClient(t, http.StatusOK, xmlHeader(), string(body))

	resp, err := client.GetPlays(t.Context(), &boardgamegeek.PlaysOptions{ID: 13, Type: boardgamegeek.PlaysTypeThing})
	require.NoError(t, err)

	assert.Empty(t, resp.Username)
	require.Len(t, resp.Plays, 2)

	play := resp.Plays[0]
	assert.EqualValues(t, 2001, play.UserID)
	assert.Equal(t, 75, play.Length)
	assert.Equal(t, "Example Game Cafe", play.Location)
	assert.Equal(t, "Mock comment with a #hashtag", play.Comments)
	require.Len(t, play.Players, 2)
	assert.Equal(t, boardgamegeek.Player{
		XMLName:       play.Players[0].XMLName,
		Username:      "examplelogger",
		UserID:        2001,
		Name:          "Alex Example",
		StartPosition: "start player",
		Color:         "Red",
		Score:         "10",
		Rating:        7,
		Win:           true,
	}, *play.Players[0])
	assert.True(t, play.Players[1].New)
	assert.Empty(t, play.Players[1].Score)

	assert.True(t, resp.Plays[1].Incomplete)
	assert.Empty(t, resp.Plays[1].Players)
}
