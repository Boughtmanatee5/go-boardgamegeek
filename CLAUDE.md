# CLAUDE.md -- Boardgamegeek Client

## Startup Rules

Before writing any code, complete these steps in order:

1. **Read this file completely.** It defines the boundaries and conventions for this project.

## Project Description
This is a go client to interface with boardgamegeek.com's api xmlapi2.
The goal is to build clean go structs and interfaces to simplify making API requests to boardgamegeek.

## Architecture
0. The module is `github.com/Boughtmanatee5/go-boardgamegeek` and the package is `boardgamegeek`, so callers import it without an alias.
   Naming follows Go conventions:
    - request params are `XxxOptions` and responses `XxxResponse`; methods return `(*XxxResponse, error)`
    - initialisms are upper case (`ID`, `CollID`)
    - enum constants are prefixed with their type (`ThingTypeBoardGame`, `HotItemTypeRPG`)
1. https://boardgamegeek.com/wiki/page/BGG_XML_API2 is the source of truth for boardgamegeek xmlapi2 don't invent fields; check them against a real response.
2. client.go contains the main struct Client.
    - Client contains all the methods that make requests to boardgamegeek
    - Build clients with functional options: `NewClient(opts ...ClientOption)`, e.g. `NewClient(WithAuthToken(token))`.
      Options are `WithAuthToken` (sends `Authorization: Bearer <token>`), `WithHTTPClient`, `WithBaseURL` and `WithUserAgent`; Client fields stay unexported.
      New client settings get a `WithXxx` option rather than an exported field
    - Every endpoint method goes through the unexported `Client.get(ctx, path, opts, v)`, which encodes the request struct with `github.com/google/go-querystring`, sends it, checks the status and unmarshals the xml.
      `get` returns the public error types in errors.go: `*AcceptedError` (202), `*RateLimitError` (429, with `RetryAfter`) and `*ErrorResponse`.
      `*ErrorResponse` covers any other status and any 200 `<error>`/`<errors>`/non-xml body.
      Endpoint methods wrap errors with context (`fmt.Errorf("thread %d: %w", ...)`) so `errors.Is`/`errors.As` still work
    - `ErrNotFound` matches every not-found shape through `ErrorResponse.Is`. `get` marks 404s and `<error>` documents.
      Endpoints mark their own quirks with `markNotFound`, so an unrelated html page or validation error is not reported as not found
3. Each boardgamegeek entity has its own file. E.g. search.go, thing.go and family_item.go contain definitions for the SearchItem, Thing and FamilyItem entities
    - forum.go covers both `/forumlist` (ForumListOptions/ForumListResponse) and `/forum` (ForumOptions/ForumResponse); thread.go covers `/thread`
    - hot.go, user.go, guild.go and collection.go cover `/hot`, `/user`, `/guild` and `/collection`
    - Shared element types live in search.go: `Name`, and `ValueAttr`, a string type whose `UnmarshalXML` reads the `value` attribute of any `<x value="..."/>` element. Use `ValueAttr` for those elements instead of a wrapper struct, and `*ValueAttr` when the element is optional
    - BGG often signals "not found" with HTTP 200:
        - `/forum` returns an HTML page (GetForum marks it with `markNotFound`)
        - `/thread` and `/guild` return `<error message="..."/>`, parsed by the shared `apiError` in errors.go
        - `/collection` returns `<errors><error><message>`, parsed by `collectionErrors`. GetCollection marks "Invalid username specified" as not found
        - `/user` is the exception: an unknown name is a 404
    - `/collection` returns HTTP 202 while BGG builds the collection. GetCollection returns an `*AcceptedError` that matches `ErrCollectionNotReady`, and leaves retrying to the caller
    - Collection status filters are `*bool` (nil = no filter, true = include, false = exclude); build them with `boardgamegeek.Bool`. Tagged `int,omitempty`, go-querystring omits nil and sends 1/0
4. Integration tests for each boardgamegeek api endpoint are written in integration_test
    - Offline unit tests live in the root package (`go test .`). query_test.go checks each request struct's encoded query string. errors_test.go stubs BGG with `httptest` through `WithBaseURL`, for paths BGG can't produce on demand, such as 429
    - `integration_test/main_test.go` reads the BGG bearer token from the `BGG_API_TOKEN` env var in `TestMain` and fails fast if it is unset; build clients with `newTestClient()`
    - Run with `BGG_API_TOKEN=<token> go test ./integration_test/...`
    - `newTestClient()` shares one `http.Client` whose `throttledTransport` spaces every request 2s apart (5s for `/collection`, which is limited more tightly) to stay under BGG's rate limit (429 after ~40 requests in a burst)

## Adding a new boardgamegeek XMLAPI2 endpoint
1. Create a file for the new endpoint
2. Define an XxxOptions struct with `url:"..."` tags for every field (go-querystring): `comma` for lists, `int,omitempty` for flags BGG takes as 1/0, `omitempty` for optional params. Add the expected encoding to query_test.go
3. Add response structs with `xml` tags taken from the API wiki. Optional sections are pointers so they come back nil when absent.
4. Add a `Client` method (receiver `c`, parameter `opts *XxxOptions`) that calls `c.get` and maps any endpoint-specific "not found" error
5. Add tests under `integration_test/`
6. Update `features.json`

## Definition of Done
1. `golangci-lint run` passes with no errors
2. All unit tests (`go test .`) and integration tests pass without errors
3. The feature appears in `feature.json` with status `"done"` and evidence (time stamp the test command was run).
4. The Architecture section of this document is updated to reflect any changes in architecture

## Refactor
`docs/refactor_plan.md` describes a phased refactor toward go-github/stripe-go conventions. Track phase status in `refactor_phases.json` using the same Definition of Done.

## Session Handoff

When resuming work, read `session-handoff.md` for context from the previous session. When finishing a session, update it with:

- What was accomplished
- What remains
- Any blockers or decisions made
- Files that were modified