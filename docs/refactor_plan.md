# Refactor Plan: Align with Standard Go API Client Libraries

Reference libraries: [google/go-github](https://github.com/google/go-github) and [stripe/stripe-go](https://github.com/stripe/stripe-go).
Phase status is tracked in [`refactor_phases.json`](../refactor_phases.json).

## Current State vs. Reference Libraries

The data layer is close to idiomatic. It has per-entity files, typed enums, pointers for optional sections,
`ctx` as the first argument, `%w` wrapping, the `Bool()` helper and the `ErrCollectionNotReady` sentinel.
The gaps are mostly in the transport layer (`client.go`), plus errors, the package layout and tests.

| Area | go-github / stripe-go | BoardGameGeekClient | Gap |
|---|---|---|---|
| Request pipeline | One `NewRequest` and one `Do(ctx, req, v)` shared by every endpoint | 10 copies of the same ~45-line build/auth/do/read/unmarshal code | Large |
| Client construction | Injectable `*http.Client`, exported `BaseURL`, `UserAgent`, `WithAuthToken`. stripe-go adds `BackendConfig` (retries, logger) | `http.Client` stored by value and not injectable, hardcoded URL constant, no User-Agent, `oauth2.Token` used only to set one static header | Large |
| Errors | Typed errors (`*ErrorResponse`, `*RateLimitError`, `*AcceptedError`) that work with `errors.As` | `fmt.Errorf` strings, so callers can't tell "not found" from a 429 or a 500 | Large |
| Query encoding | Struct tags encoded by `go-querystring` (go-github) or a `form` encoder (stripe-go) | A hand-written `toQueryString` per request, mostly `fmt.Sprintf`. `SearchRequest` only escapes spaces, so `&` in a query breaks it | Medium (bug) |
| API surface | Grouped services: `client.Repositories.Get`, `sc.V1Customers.Retrieve` | All methods sit directly on `Client` | Medium |
| Return values | `(*T, *Response, error)`, where `*Response` exposes the status, headers and rate info | `(*T, error)` | Medium |
| Pagination | Shared `ListOptions` and `NextPage` (go-github), iterators (stripe-go) | `Page`/`PageSize` repeated per request, no helper | Medium |
| Rate limiting / retries | `MaxNetworkRetries` with backoff, rate-limit detection | None. The tests throttle with `time.Sleep` (~4 minute suite) | Medium |
| Testing | Unit tests against `httptest.Server` with fixtures, plus optional live tests | Live integration tests only | Large |
| Naming / idiom | `ID`, `ThingTypeBoardGame`, short receivers, godoc everywhere | `Id`, ALL-CAPS constants, `bggClient` receiver, little godoc, no `doc.go` or examples | Small (breaking) |
| Package path | Import path matches the package name | Module `.../BoardGameGeekClient`, package `boardgamegeek`, so every caller needs an import alias | Small (breaking) |

Out of scope for this refactor but worth noting: `/plays` is not implemented, and `Thing` is missing many wiki fields.

## Phases

Phases 1–4 are internal or additive. The breaking changes are grouped into Phase 5 so they ship in one release.

### Phase 1: Shared request pipeline (no API change)
- Add one unexported helper on `Client` that builds the URL, creates the request, sets auth,
  sends it, checks the status, reads the body and unmarshals it.
- Detect the BGG "not found" quirks in that helper and return unexported sentinel/typed errors
  for the non-XML response, the `<error>` root, the `<errors>` root and HTTP 202.
  Each endpoint method maps these to its current error message, so callers see no change.
- Each endpoint method shrinks to a few lines.
- Fix the `boadgamegeekURL` and `recieved` typos.

### Phase 2: Client construction
- `NewClient(httpClient *http.Client) *Client` (nil uses a default), with an exported `BaseURL *url.URL`
  and `UserAgent`, plus `WithAuthToken(token string) *Client`.
- Optionally drop the `oauth2` dependency, since BGG's token is a static bearer token.
- Integration tests inject a rate-limited `RoundTripper` instead of sleeping in `newTestClient()`.

### Phase 3: Typed errors
- Add `ErrorResponse{Response *http.Response; Message string}`, `RateLimitError{RetryAfter time.Duration}`
  and `AcceptedError`, plus an `ErrNotFound` sentinel that every not-found shape matches with `errors.Is`.
- Keep `ErrCollectionNotReady` working through `errors.Is` on `AcceptedError`.

### Phase 4: Declarative query encoding
- Replace each `toQueryString` with `url:"..."` struct tags encoded by `github.com/google/go-querystring`
  (`comma` for lists, `int` for 0/1 bools, `omitempty`).
- Add a custom encoder for the collection `*bool` filters (nil is omitted, true/false become 1/0).
- This fixes the search escaping bug.

### Phase 5: Services and naming (breaking)
- go-github-style services, where each is `type XService service` and all share one `*Client`:
  - `Things.Get` / `Things.Search` / `Things.Hot`
  - `Families.Get`
  - `Forums.List` / `Forums.Get` / `Forums.GetThread`
  - `Users.Get`
  - `Guilds.Get`
  - `Collections.Get`
- Rename the request types from `XxxRequest` to `XxxOptions`, and return `(*T, *Response, error)`.
- Rename `Id` to `ID` and the constants to `ThingTypeBoardGame` style, and use short receivers.
- Fix the module path/package name mismatch (rename the module, or move the code into a subdirectory).

### Phase 6: Pagination and retries
- Embed a shared `ListOptions{Page, PageSize}`.
- Add `iter.Seq2[*T, error]` iterators for the paged lists (thing comments, forum threads, guild members, user buddies).
- Add opt-in retries for 202/429 with backoff, configured on `Client`.

### Phase 7: Tests and docs
- Add unit tests in the root package against `httptest.NewServer`, with real BGG responses saved in `testdata/*.xml`.
- Put `integration_test/` behind a `//go:build integration` tag.
- Add `doc.go`, godoc on exported identifiers, `example_test.go` and a real README.

## Per-Phase Definition of Done
Same as `CLAUDE.md`:
1. `golangci-lint run` passes.
2. The integration tests pass.
3. The phase is marked `"done"` in `refactor_phases.json` with evidence (timestamp and test command).
4. `CLAUDE.md` Architecture and "Adding a new endpoint" sections reflect the change.
