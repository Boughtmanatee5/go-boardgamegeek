package boardgamegeek

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ErrNotFound matches, with errors.Is, every way BGG reports a missing entity:
// a 404, a 200 <error> document, /forum's html page and /collection's invalid username error
var ErrNotFound = errors.New("boardgamegeek: not found")

// ErrorResponse is returned when BGG rejects a request, either with a non-200 status
// or with a 200 error document in place of the expected response
type ErrorResponse struct {
	// Response is the http response; its body has already been read and closed
	Response *http.Response
	// Message is BGG's error message, or the http status when BGG sent none
	Message string

	notFound bool
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, e.Response.Request.URL, e.Response.StatusCode, e.Message)
}

// Is reports whether target is ErrNotFound and the response means the entity does not exist
func (e *ErrorResponse) Is(target error) bool {
	return target == ErrNotFound && e.notFound
}

// RateLimitError is returned for HTTP 429 Too Many Requests
type RateLimitError struct {
	Response *http.Response
	// RetryAfter is how long BGG asked the caller to wait, or 0 if it did not say
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("%s %s: rate limited by boardgamegeek, retry after %s", e.Response.Request.Method, e.Response.Request.URL, e.RetryAfter)
	}
	return fmt.Sprintf("%s %s: rate limited by boardgamegeek", e.Response.Request.Method, e.Response.Request.URL)
}

// AcceptedError is returned for HTTP 202, which BGG sends while it queues a request.
// Retry the same request after a few seconds.
type AcceptedError struct {
	Response *http.Response
}

func (e *AcceptedError) Error() string {
	return fmt.Sprintf("%s %s: request queued by boardgamegeek, retry later", e.Response.Request.Method, e.Response.Request.URL)
}

// Is lets errors.Is(err, ErrCollectionNotReady) keep matching a queued collection request
func (e *AcceptedError) Is(target error) bool {
	return target == ErrCollectionNotReady
}

// apiError is returned with HTTP 200 in place of the expected document when an entity
// does not exist, e.g. <error message='Thread Not Found' /> from /thread and /guild
type apiError struct {
	XMLName xml.Name `xml:"error"`
	Message string   `xml:"message,attr"`
}

// parseRetryAfter reads a Retry-After header given in seconds or as an http date
func parseRetryAfter(header string) time.Duration {
	if header == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(header); err == nil {
		return time.Duration(seconds) * time.Second
	}
	if date, err := http.ParseTime(header); err == nil {
		return max(time.Until(date), 0)
	}
	return 0
}
