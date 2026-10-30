# go-boardgamegeek

A Go client for the [BoardGameGeek XML API2](https://boardgamegeek.com/wiki/page/BGG_XML_API2).

```sh
go get github.com/Boughtmanatee5/go-boardgamegeek
```

Requires Go 1.26 (see `go.mod`).

## Usage

BGG requires an API token for every request; see the [API wiki](https://boardgamegeek.com/wiki/page/BGG_XML_API2) for how to get one.

```go
import "github.com/Boughtmanatee5/go-boardgamegeek"

client := boardgamegeek.NewClient(
    boardgamegeek.WithAuthToken(os.Getenv("BGG_API_TOKEN")),
    boardgamegeek.WithMaxRetries(5), // retry 202 (queued) and 429 (rate limited) responses
)

resp, err := client.GetThings(ctx, &boardgamegeek.ThingsOptions{IDs: []int64{13}, Stats: true})
```

### Client options

| Option | Default |
|---|---|
| `WithAuthToken(token)` | none, but BGG rejects unauthenticated requests |
| `WithHTTPClient(*http.Client)` | a new `http.Client` |
| `WithMaxRetries(n)` | 0 (no retries) |
| `WithBaseURL(*url.URL)` | `https://boardgamegeek.com/xmlapi2/` |
| `WithUserAgent(string)` | `go-boardgamegeek` |

### Endpoints

| BGG endpoint | Method | Iterator |
|---|---|---|
| `/search` | `Search` | |
| `/thing` | `GetThings` | `AllThingComments` |
| `/family` | `GetFamilyItems` | |
| `/forumlist` | `GetForumList` | |
| `/forum` | `GetForum` | `AllForumThreads` |
| `/thread` | `GetThread` | |
| `/hot` | `GetHotItems` | |
| `/user` | `GetUser` | `AllUserBuddies`, `AllUserGuilds` |
| `/guild` | `GetGuild` | `AllGuildMembers` |
| `/collection` | `GetCollection` | |
| `/plays` | `GetPlays` | `AllPlays` |

The iterators return an `iter.Seq2` and request pages as the loop needs them:

```go
for thread, err := range client.AllForumThreads(ctx, &boardgamegeek.ForumOptions{ID: 1508}) {
    if err != nil {
        return err
    }
    fmt.Println(thread.Subject)
}
```

### Errors

BGG often reports a missing entity with HTTP 200 and an error document. Every form of "not found" matches `ErrNotFound`:

```go
_, err := client.GetUser(ctx, &boardgamegeek.UserOptions{Name: "nobody"})
if errors.Is(err, boardgamegeek.ErrNotFound) { ... }
```

Use `errors.As` for `*ErrorResponse` (the HTTP response and BGG's message), `*RateLimitError` (429, with `RetryAfter`)
or `*AcceptedError` (202, sent while BGG builds a collection).

### Rate limits

BGG returns 429 after a burst of roughly 40 requests, and `/collection` is limited more tightly.
Space requests out, for example with a throttling `http.RoundTripper` passed through `WithHTTPClient`, and set `WithMaxRetries`.

## Development

The integration tests call the live API and read your token from `BGG_API_TOKEN`. Copy `.env.example` to `.env` (git-ignored) to keep it locally.

```sh
go test ./...                                                       # offline unit tests
BGG_API_TOKEN=<token> go test -tags integration ./integration_test/... # live tests, ~4 minutes
golangci-lint run
```
