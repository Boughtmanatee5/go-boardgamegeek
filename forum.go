package boardgamegeek

import (
	"encoding/xml"
)

type ForumListType string

const (
	ForumListTypeThing  ForumListType = "thing"
	ForumListTypeFamily ForumListType = "family"
)

type ForumListRequest struct {
	ID   int64         `url:"id"`
	Type ForumListType `url:"type"`
}

type ForumSummary struct {
	XMLName      xml.Name `xml:"forum"`
	Id           int64    `xml:"id,attr"`
	GroupId      int64    `xml:"groupid,attr"`
	Title        string   `xml:"title,attr"`
	NoPosting    bool     `xml:"noposting,attr"`
	Description  string   `xml:"description,attr"`
	NumThreads   int      `xml:"numthreads,attr"`
	NumPosts     int      `xml:"numposts,attr"`
	LastPostDate string   `xml:"lastpostdate,attr"`
}

type ForumListResponse struct {
	XMLName    xml.Name        `xml:"forums"`
	Type       string          `xml:"type,attr"`
	Id         int64           `xml:"id,attr"`
	TermsOfUse string          `xml:"termsofuse,attr"`
	Forums     []*ForumSummary `xml:"forum"`
}

type ForumRequest struct {
	ID   int64 `url:"id"`
	Page int   `url:"page,omitempty"`
}

type ThreadSummary struct {
	XMLName      xml.Name `xml:"thread"`
	Id           int64    `xml:"id,attr"`
	Subject      string   `xml:"subject,attr"`
	Author       string   `xml:"author,attr"`
	NumArticles  int      `xml:"numarticles,attr"`
	PostDate     string   `xml:"postdate,attr"`
	LastPostDate string   `xml:"lastpostdate,attr"`
}

type ForumResponse struct {
	XMLName    xml.Name `xml:"forum"`
	Id         int64    `xml:"id,attr"`
	Title      string   `xml:"title,attr"`
	NumThreads int      `xml:"numthreads,attr"`
	NumPosts   int      `xml:"numposts,attr"`
	// BGG always returns the Unix epoch here; use ForumSummary.LastPostDate from GetForumList instead.
	LastPostDate string           `xml:"lastpostdate,attr"`
	NoPosting    bool             `xml:"noposting,attr"`
	TermsOfUse   string           `xml:"termsofuse,attr"`
	Threads      []*ThreadSummary `xml:"threads>thread"`
}
