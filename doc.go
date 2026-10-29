// Package boardgamegeek is a client for the BoardGameGeek XML API2
// (https://boardgamegeek.com/wiki/page/BGG_XML_API2).
//
// Create a Client with NewClient and functional options. BGG requires an API token:
//
//	client := boardgamegeek.NewClient(
//		boardgamegeek.WithAuthToken(os.Getenv("BGG_API_TOKEN")),
//		boardgamegeek.WithMaxRetries(5),
//	)
//
//	resp, err := client.Search(ctx, &boardgamegeek.SearchOptions{Query: "Catan"})
//
// Each endpoint has a Client method that takes an XxxOptions struct and returns an XxxResponse.
// Paged lists (forum threads, guild members, user buddies and guilds, thing comments) also have
// AllXxx methods that return an iter.Seq2 and request pages as the loop needs them.
//
// # Errors
//
// BGG often reports a missing entity with HTTP 200, using an error document or an HTML page.
// Every such response, and a 404, matches ErrNotFound:
//
//	if errors.Is(err, boardgamegeek.ErrNotFound) { ... }
//
// Use errors.As to get an *ErrorResponse, a *RateLimitError (HTTP 429) or an *AcceptedError
// (HTTP 202, sent while BGG builds a collection). WithMaxRetries retries the last two for you.
//
// # Rate limits
//
// BGG rate limits clients and returns 429 after a burst of roughly 40 requests. Space requests out,
// for example with a throttling http.RoundTripper passed through WithHTTPClient, or set WithMaxRetries.
package boardgamegeek
