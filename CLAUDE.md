# CLAUDE.md -- Boardgamegeek Client

## Startup Rules

Before writing any code, complete these steps in order:

1. **Read this file completely.** It defines the boundaries and conventions for this project.

## Project Description
This is a go client to interface with boardgamegeek.com's api xmlapi2.
The goal is to build clean go structs and interfaces to simplify making API requests to boardgamegeek.

## Architecture
0. The module targets Go 1.26 (go.mod `go 1.26`) so older toolchains can use it; prefer stdlib features up to 1.26 (`errors.AsType`, `new(expr)`). The module is `github.com/Boughtmanatee5/go-boardgamegeek` and the package is `boardgamegeek`, so callers import it without an alias.
   Naming follows Go conventions:
    - request params are `XxxOptions` and responses `XxxResponse`; methods return `(*XxxResponse, error)`
    - initialisms are upper case (`ID`, `CollID`)
    - enum constants are prefixed with their type (`ThingTypeBoardGame`, `HotItemTypeRPG`)
1. https://boardgamegeek.com/wiki/page/BGG_XML_API2 is the source of truth for boardgamegeek xmlapi2 don't invent fields; check them against a real response.
2. client.go contains the main struct Client.
    - Client contains all the methods that make requests to boardgamegeek
    - Build clients with functional options: `NewClient(opts ...ClientOption)`, e.g. `NewClient(WithAuthToken(token))`.
      Options are `WithAuthToken` (sends `Authorization: Bearer <token>`), `WithHTTPClient`, `WithBaseURL`, `WithUserAgent` and `WithMaxRetries`; Client fields stay unexported.
      New client settings get a `WithXxx` option rather than an exported field
    - Every endpoint method goes through the unexported `Client.get(ctx, path, opts, v)`, which encodes the request struct with `github.com/google/go-querystring`, sends it, checks the status and unmarshals the xml.
      `get` returns the public error types in errors.go: `*AcceptedError` (202), `*RateLimitError` (429, with `RetryAfter`) and `*ErrorResponse`.
      `*ErrorResponse` covers any other status and any 200 `<error>`/`<errors>`/messagebox/non-xml body.
      Endpoint methods wrap errors with context (`fmt.Errorf("thread %d: %w", ...)`) so `errors.Is`/`errors.As` still work
    - `ErrNotFound` matches every not-found shape through `ErrorResponse.Is`. `get` marks 404s and `<error>` documents.
      Endpoints mark their own quirks with `markNotFound`, so an unrelated html page or validation error is not reported as not found
    - `get` encodes the query and loops over `do`, which sends a single request. With `WithMaxRetries(n)`, a 202 or 429 is retried up to n times.
      The wait is `Retry-After` when sent, else backoff from 2s (`retryWait`) doubling to 1 minute, and the wait stops when ctx is done. The default is no retries
    - pagination.go holds the generic `paginate` helper and the `AllXxx` iterators (`AllForumThreads`, `AllGuildMembers`, `AllUserBuddies`, `AllUserGuilds`, `AllThingComments`, `AllPlays`), which return `iter.Seq2[*T, error]`.
      `paginate` stops on an empty or short page, on reaching BGG's reported total, or on the first error. Each page copies opts, so the caller's options are never mutated.
      There is no shared `ListOptions`: only `/thing` accepts `pagesize`, so each options struct keeps its own `Page`
    - doc.go is the package overview. Every exported identifier needs a doc comment (enforced by revive's `exported` rule in .golangci.yml)
3. Each boardgamegeek entity has its own file. E.g. search.go, thing.go and family_item.go contain definitions for the SearchItem, Thing and FamilyItem entities
    - forum.go covers both `/forumlist` (ForumListOptions/ForumListResponse) and `/forum` (ForumOptions/ForumResponse); thread.go covers `/thread`
    - hot.go, user.go, guild.go, collection.go and plays.go cover `/hot`, `/user`, `/guild`, `/collection` and `/plays`
    - `ThingType` (thing.go) filters `/thing` and `/collection`. `SearchType` (search.go) filters `/search`: it repeats every thing type and adds types that are not things, such as `boardgamedesigner`
      (designers have their own ids; `/thing?id=2&type=boardgamedesigner` returns no items, while `/search` finds designer 2, Reiner Knizia)
    - Shared element types live in search.go: `Name`, `Link` (thing, version and family links, with `Inbound`), and the value-attribute types for `<x value="..."/>` elements:
      `ValueAttr` (string), `IntValue` and `FloatValue`, whose `UnmarshalXML` reads the `value` attribute (an empty value is 0; a non-number is an error).
      Use them instead of a wrapper struct, pick `IntValue`/`FloatValue` for counts, years, times and averages, keep `ValueAttr` for text and for values BGG sends unparsed (`"Not Ranked"`, `"2019-00-00"`), and use pointers when the element is optional
    - thing.go: `Thing` covers every thing type; elements a type lacks are nil (player counts and polls for board games, `ReleaseDate` for video games, `SeriesCode` for RPG items, `DatePublished`/`IssueIndex` for RPG issues).
      Video game "versions" are characters with only a canonical name
    - BGG often signals "not found" with HTTP 200:
        - `/forum` returns an HTML page (GetForum marks it with `markNotFound`)
        - `/thread` and `/guild` return `<error message="..."/>`, parsed by the shared `apiError` in errors.go
        - `/collection` returns `<errors><error><message>`, parsed by `collectionErrors`. GetCollection marks "Invalid username specified" as not found
        - `/plays` returns an xml declaration followed by `<div class='messagebox error'>Invalid object or user</div>`, parsed by `messageBox` in plays.go. GetPlays marks it as not found.
          With neither username nor id, `/plays` returns `<error message='Not Found'/>`, so GetPlays rejects that request itself rather than report it as not found
        - `/user` is the exception: an unknown name is a 404
    - `/collection` returns HTTP 202 while BGG builds the collection. GetCollection returns an `*AcceptedError` that matches `ErrCollectionNotReady`; callers retry themselves or set `WithMaxRetries`
    - Collection status filters are `*bool` (nil = no filter, true = include, false = exclude); build them with `new(true)` / `new(false)` (Go 1.26). Tagged `int,omitempty`, go-querystring omits nil and sends 1/0
4. Tests
    - Offline unit tests live in the root package and run with `go test ./...` (no token needed):
        - endpoints_test.go parses real BGG responses saved in `testdata/<endpoint>.xml` (each test notes the request it was captured with). Re-capture a fixture when a response struct changes.
          Don't commit fixtures holding other users' records (e.g. `/plays` by thing); hand-write a mock with fictional data instead, like `testdata/plays_players.xml`
        - query_test.go checks each options struct's encoded query string
        - errors_test.go and retry_test.go stub BGG with `httptest` through `WithBaseURL`, for paths BGG can't produce on demand, such as 429
        - pagination_test.go covers `paginate`'s stop conditions
        - example_test.go holds godoc examples. They call the live API, so they have no `// Output:` and are compiled, not run
    - Integration tests for each boardgamegeek api endpoint are written in integration_test, behind the `//go:build integration` tag. .golangci.yml sets that build tag so they are still linted
    - `integration_test/main_test.go` reads the BGG bearer token from the `BGG_API_TOKEN` env var in `TestMain` and fails fast if it is unset; build clients with `newTestClient(opts...)`, which takes extra client options
    - Run with `BGG_API_TOKEN=<token> go test -tags integration ./integration_test/...`
    - CI (.github/workflows/ci.yml) runs golangci-lint, `go mod tidy -diff` and `go test -race ./...` on Go 1.26 and stable, for pushes to main and pull requests.
      Never add the integration tests to CI: repeated live runs could get the BGG token rate limited or revoked. Keep the lint action's golangci-lint version in step with the local one
    - Collection tests use `WithMaxRetries` to wait out BGG's 202s, which exercises the retry path live
    - `newTestClient()` shares one `http.Client` whose `throttledTransport` spaces every request 2s apart (5s for `/collection`, which is limited more tightly) to stay under BGG's rate limit (429 after ~40 requests in a burst)

## Adding a new boardgamegeek XMLAPI2 endpoint
1. Create a file for the new endpoint
2. Define an XxxOptions struct with `url:"..."` tags for every field (go-querystring): `comma` for lists, `int,omitempty` for flags BGG takes as 1/0, `omitempty` for optional params. Add the expected encoding to query_test.go
3. Add response structs with `xml` tags taken from the API wiki. Optional sections are pointers so they come back nil when absent.
4. Add a `Client` method (receiver `c`, parameter `opts *XxxOptions`) that calls `c.get` and maps any endpoint-specific "not found" error
5. Capture a real response into `testdata/<endpoint>.xml` and add a parsing test to endpoints_test.go. Add integration tests under `integration_test/` (with the `integration` build tag). Add an `AllXxx` iterator in pagination.go if the endpoint is paged
6. Update `features.json`

## Definition of Done
1. `golangci-lint run` passes with no errors
2. All unit tests (`go test ./...`) and integration tests (`go test -tags integration ./integration_test/...`) pass without errors
3. The feature appears in `features.json` with status `"done"` and evidence (time stamp the test command was run).
4. The Architecture section of this document is updated to reflect any changes in architecture

## Design decisions
These were settled during the refactor toward go-github/stripe-go conventions; don't re-propose them without a new reason.
- Endpoint methods live directly on `Client`, not in go-github-style services. BGG has only about 10 endpoints, so grouping them adds indirection without much benefit
- Methods return `(*T, error)`, not `(*T, *Response, error)`. BGG returns paging in the xml body and sends no rate-limit headers, and the typed errors already carry the `*http.Response`
- The module is `github.com/Boughtmanatee5/go-boardgamegeek` and the package is `boardgamegeek` (the go-github convention)
- There is no shared `ListOptions{Page, PageSize}`. Only `/thing` accepts `pagesize`, so a shared struct would add a field BGG ignores to the forum, guild and user options

## Session Handoff
`session-handoff.md` is an optional, local notes file (it is git-ignored, so it is not shared with other contributors).
If it exists, read it when resuming work for context from the previous session. When finishing a session, update it with:

- What was accomplished
- What remains
- Any blockers or decisions made
- Files that were modified
