# go-boardgamegeek

[![CI](https://github.com/Boughtmanatee5/go-boardgamegeek/actions/workflows/ci.yml/badge.svg)](https://github.com/Boughtmanatee5/go-boardgamegeek/actions/workflows/ci.yml)

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
CI runs lint and the unit tests only; run the integration tests yourself before sending a change to an endpoint.

```sh
go test ./...                                                       # offline unit tests
BGG_API_TOKEN=<token> go test -tags integration ./integration_test/... # live tests, ~4 minutes
golangci-lint run
```

## Contributing

Contributions are welcome. Please open an issue before starting on anything large.

### Issues

Open issues from one of the templates so they are labelled:

- **Bug report** (`bug`): something is broken or a response fails to parse. Include the request (endpoint and parameters), the library version and what you expected.
- **Feature request** (`enhancement`): something new, such as an endpoint, field or client option.
- **Improvement** (`enhancement`): a change to existing behaviour, docs or performance.

### Pull requests

PRs must be submitted from a fork of this repository; branches are not accepted on the main repository itself.

1. Fork the repository and create a branch from `main`.
2. Make your change, following the conventions in [CLAUDE.md](CLAUDE.md): `XxxOptions`/`XxxResponse` naming, upper-case initialisms and a doc comment on every exported identifier.
3. Run `golangci-lint run` and `go test ./...`. If you changed an endpoint, run the integration tests as well (see [Development](#development)).
4. Open a PR against `main`, link the issue it addresses and say which tests you ran.
