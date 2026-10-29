package boardgamegeek

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type ThreadRequest struct {
	ID           int64
	MinArticleID int64
	// MinArticleDate is YYYY-MM-DD or YYYY-MM-DD HH:MM:SS
	MinArticleDate string
	Count          int
}

func (r *ThreadRequest) toQueryString() string {
	query := fmt.Sprintf("id=%d", r.ID)

	if r.MinArticleID > 0 {
		query = fmt.Sprintf("%s&minarticleid=%d", query, r.MinArticleID)
	}

	if r.MinArticleDate != "" {
		query = fmt.Sprintf("%s&minarticledate=%s", query, url.QueryEscape(r.MinArticleDate))
	}

	if r.Count > 0 {
		query = fmt.Sprintf("%s&count=%d", query, r.Count)
	}

	return query
}

type Article struct {
	XMLName  xml.Name `xml:"article"`
	Id       int64    `xml:"id,attr"`
	Username string   `xml:"username,attr"`
	Link     string   `xml:"link,attr"`
	PostDate string   `xml:"postdate,attr"`
	EditDate string   `xml:"editdate,attr"`
	NumEdits int      `xml:"numedits,attr"`
	Subject  string   `xml:"subject"`
	Body     string   `xml:"body"`
}

type ThreadResponse struct {
	XMLName     xml.Name   `xml:"thread"`
	Id          int64      `xml:"id,attr"`
	NumArticles int        `xml:"numarticles,attr"`
	Link        string     `xml:"link,attr"`
	TermsOfUse  string     `xml:"termsofuse,attr"`
	Subject     string     `xml:"subject"`
	Articles    []*Article `xml:"articles>article"`
}
