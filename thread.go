package boardgamegeek

import (
	"encoding/xml"
)

// ThreadOptions are the parameters for GetThread
type ThreadOptions struct {
	ID           int64 `url:"id"`
	MinArticleID int64 `url:"minarticleid,omitempty"`
	// MinArticleDate is YYYY-MM-DD or YYYY-MM-DD HH:MM:SS
	MinArticleDate string `url:"minarticledate,omitempty"`
	Count          int    `url:"count,omitempty"`
}

// Article is one post in a thread
type Article struct {
	XMLName  xml.Name `xml:"article"`
	ID       int64    `xml:"id,attr"`
	Username string   `xml:"username,attr"`
	Link     string   `xml:"link,attr"`
	PostDate string   `xml:"postdate,attr"`
	EditDate string   `xml:"editdate,attr"`
	NumEdits int      `xml:"numedits,attr"`
	Subject  string   `xml:"subject"`
	Body     string   `xml:"body"`
}

// ThreadResponse is the response from GetThread
type ThreadResponse struct {
	XMLName     xml.Name   `xml:"thread"`
	ID          int64      `xml:"id,attr"`
	NumArticles int        `xml:"numarticles,attr"`
	Link        string     `xml:"link,attr"`
	TermsOfUse  string     `xml:"termsofuse,attr"`
	Subject     string     `xml:"subject"`
	Articles    []*Article `xml:"articles>article"`
}
