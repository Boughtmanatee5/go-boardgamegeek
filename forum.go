package boardgamegeek

import (
	"encoding/xml"
)

// ForumListType is the kind of entity whose forums GetForumList returns
type ForumListType string

// Entity types accepted by ForumListOptions.Type
const (
	ForumListTypeThing  ForumListType = "thing"
	ForumListTypeFamily ForumListType = "family"
)

// ForumListOptions are the parameters for GetForumList
type ForumListOptions struct {
	ID   int64         `url:"id"`
	Type ForumListType `url:"type"`
}

// ForumSummary describes one forum in a ForumListResponse
type ForumSummary struct {
	XMLName      xml.Name `xml:"forum"`
	ID           int64    `xml:"id,attr"`
	GroupID      int64    `xml:"groupid,attr"`
	Title        string   `xml:"title,attr"`
	NoPosting    bool     `xml:"noposting,attr"`
	Description  string   `xml:"description,attr"`
	NumThreads   int      `xml:"numthreads,attr"`
	NumPosts     int      `xml:"numposts,attr"`
	LastPostDate string   `xml:"lastpostdate,attr"`
}

// ForumListResponse is the response from GetForumList
type ForumListResponse struct {
	XMLName    xml.Name        `xml:"forums"`
	Type       string          `xml:"type,attr"`
	ID         int64           `xml:"id,attr"`
	TermsOfUse string          `xml:"termsofuse,attr"`
	Forums     []*ForumSummary `xml:"forum"`
}

// ForumOptions are the parameters for GetForum and AllForumThreads
type ForumOptions struct {
	ID   int64 `url:"id"`
	Page int   `url:"page,omitempty"`
}

// ThreadSummary describes one thread in a forum; use GetThread for its articles
type ThreadSummary struct {
	XMLName      xml.Name `xml:"thread"`
	ID           int64    `xml:"id,attr"`
	Subject      string   `xml:"subject,attr"`
	Author       string   `xml:"author,attr"`
	NumArticles  int      `xml:"numarticles,attr"`
	PostDate     string   `xml:"postdate,attr"`
	LastPostDate string   `xml:"lastpostdate,attr"`
}

// ForumResponse is the response from GetForum, one page of a forum's threads
type ForumResponse struct {
	XMLName    xml.Name `xml:"forum"`
	ID         int64    `xml:"id,attr"`
	Title      string   `xml:"title,attr"`
	NumThreads int      `xml:"numthreads,attr"`
	NumPosts   int      `xml:"numposts,attr"`
	// BGG always returns the Unix epoch here; use ForumSummary.LastPostDate from GetForumList instead.
	LastPostDate string           `xml:"lastpostdate,attr"`
	NoPosting    bool             `xml:"noposting,attr"`
	TermsOfUse   string           `xml:"termsofuse,attr"`
	Threads      []*ThreadSummary `xml:"threads>thread"`
}
