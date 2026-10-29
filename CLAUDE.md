# CLAUDE.md -- Boardgamegeek Client

## Startup Rules

Before writing any code, complete these steps in order:

1. **Read this file completely.** It defines the boundaries and conventions for this project.

## Project Description
This is a go client to interface with boardgamegeek.com's api xmlapi2.
The goal is to build clean go structs and interfaces to simplify making API requests to boardgamegeek.

## Architecture
1. https://boardgamegeek.com/wiki/page/BGG_XML_API2 is the source of truth for boardgamegeek xmlapi2 don't invent fields; check them against a real response.
2. client.go contains the main struct Client.
    - Client contains all the methods that make requests to boardgamegeek
3. Each boardgamegeek entity has its own file. E.g. search.go, thing.go and family_item.go contain definitions for the SearchItem, Thing and FamilyItem entities
4. Integration tests for each boardgamegeek api endpoint are written in integration_test
    - `integration_test/main_test.go` reads the BGG bearer token from the `BGG_API_TOKEN` env var in `TestMain` and fails fast if it is unset; build clients with `newTestClient()`
    - Run with `BGG_API_TOKEN=<token> go test ./integration_test/...`

## Adding a new boardgamegeek XMLAPI2 endpoint
1. Create a file for the new endpoint
2. Define XxxRequest struct and a toQuerString method
3. Add response structs with `xml` tags taken from the API wiki. Optional sections are pointers so they come back nil when absent.
4. Add a `Client` method that wraps errors with `%w` and checks the response status code
5. Add tests under `integration_test/`
6. Update `features.json`

## Definition of Done
1. `golangci-lint run` passes with no errors
2. All integration tests pass without errors
3. The feature appears in `feature.json` with status `"done"` and evidence (time stamp the test command was run).
4. The Architecture section of this document is updated to reflect any changes in architecture

## Session Handoff

When resuming work, read `session-handoff.md` for context from the previous session. When finishing a session, update it with:

- What was accomplished
- What remains
- Any blockers or decisions made
- Files that were modified