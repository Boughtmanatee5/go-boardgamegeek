package boardgamegeek

import (
	"encoding/xml"
	"errors"
)

// ErrCollectionNotReady matches, with errors.Is, the *AcceptedError returned when BGG has queued
// the collection request (HTTP 202). Retry the same request after a few seconds.
var ErrCollectionNotReady = errors.New("collection request queued by boardgamegeek, retry later")

// Bool returns a pointer to b, for the optional filters on CollectionOptions.
func Bool(b bool) *bool {
	return &b
}

type CollectionOptions struct {
	// Username is required
	Username       string    `url:"username"`
	Version        bool      `url:"version,int,omitempty"`
	Subtype        ThingType `url:"subtype,omitempty"`
	ExcludeSubtype ThingType `url:"excludesubtype,omitempty"`
	IDs            []int64   `url:"id,comma,omitempty"`
	// Brief returns only the name and status of each item
	Brief bool `url:"brief,int,omitempty"`
	Stats bool `url:"stats,int,omitempty"`
	// Status filters: nil means no filter, true includes only matching items, false excludes them
	Own        *bool `url:"own,int,omitempty"`
	Rated      *bool `url:"rated,int,omitempty"`
	Played     *bool `url:"played,int,omitempty"`
	Comment    *bool `url:"comment,int,omitempty"`
	Trade      *bool `url:"trade,int,omitempty"`
	Want       *bool `url:"want,int,omitempty"`
	Wishlist   *bool `url:"wishlist,int,omitempty"`
	Preordered *bool `url:"preordered,int,omitempty"`
	WantToPlay *bool `url:"wanttoplay,int,omitempty"`
	WantToBuy  *bool `url:"wanttobuy,int,omitempty"`
	PrevOwned  *bool `url:"prevowned,int,omitempty"`
	HasParts   *bool `url:"hasparts,int,omitempty"`
	WantParts  *bool `url:"wantparts,int,omitempty"`
	// WishlistPriority filters on wishlist priority 1-5, 0 means no filter
	WishlistPriority int `url:"wishlistpriority,omitempty"`
	// Rating filters are 1-10, 0 means no filter
	MinRating    float64 `url:"minrating,omitempty"`
	Rating       float64 `url:"rating,omitempty"`
	MinBGGRating float64 `url:"minbggrating,omitempty"`
	BGGRating    float64 `url:"bggrating,omitempty"`
	// Play count filters, 0 means no filter
	MinPlays int `url:"minplays,omitempty"`
	MaxPlays int `url:"maxplays,omitempty"`
	// ShowPrivate only returns private info for the authenticated user's own collection
	ShowPrivate bool  `url:"showprivate,int,omitempty"`
	CollID      int64 `url:"collid,omitempty"`
	// ModifiedSince is YY-MM-DD or YY-MM-DD HH:MM:SS
	ModifiedSince string `url:"modifiedsince,omitempty"`
}

// CollectionName differs from Name: the value is the element text rather than an attribute
type CollectionName struct {
	XMLName   xml.Name `xml:"name"`
	SortIndex int      `xml:"sortindex,attr"`
	Value     string   `xml:",chardata"`
}

type CollectionStatus struct {
	XMLName          xml.Name `xml:"status"`
	Own              bool     `xml:"own,attr"`
	PrevOwned        bool     `xml:"prevowned,attr"`
	ForTrade         bool     `xml:"fortrade,attr"`
	Want             bool     `xml:"want,attr"`
	WantToPlay       bool     `xml:"wanttoplay,attr"`
	WantToBuy        bool     `xml:"wanttobuy,attr"`
	Wishlist         bool     `xml:"wishlist,attr"`
	WishlistPriority int      `xml:"wishlistpriority,attr"`
	Preordered       bool     `xml:"preordered,attr"`
	LastModified     string   `xml:"lastmodified,attr"`
}

type CollectionRank struct {
	XMLName      xml.Name `xml:"rank"`
	Type         string   `xml:"type,attr"`
	ID           int64    `xml:"id,attr"`
	Name         string   `xml:"name,attr"`
	FriendlyName string   `xml:"friendlyname,attr"`
	// Value and BayesAverage are "Not Ranked" for unranked items
	Value        string `xml:"value,attr"`
	BayesAverage string `xml:"bayesaverage,attr"`
}

type CollectionRating struct {
	XMLName xml.Name `xml:"rating"`
	// Value is the user's rating, "N/A" when unrated
	Value        string            `xml:"value,attr"`
	UsersRated   ValueAttr         `xml:"usersrated"`
	Average      ValueAttr         `xml:"average"`
	BayesAverage ValueAttr         `xml:"bayesaverage"`
	StdDev       ValueAttr         `xml:"stddev"`
	Median       ValueAttr         `xml:"median"`
	Ranks        []*CollectionRank `xml:"ranks>rank"`
}

type CollectionStats struct {
	XMLName     xml.Name          `xml:"stats"`
	MinPlayers  int               `xml:"minplayers,attr"`
	MaxPlayers  int               `xml:"maxplayers,attr"`
	MinPlayTime int               `xml:"minplaytime,attr"`
	MaxPlayTime int               `xml:"maxplaytime,attr"`
	PlayingTime int               `xml:"playingtime,attr"`
	NumOwned    int               `xml:"numowned,attr"`
	Rating      *CollectionRating `xml:"rating"`
}

type CollectionItem struct {
	XMLName         xml.Name          `xml:"item"`
	ObjectType      string            `xml:"objecttype,attr"`
	ObjectID        int64             `xml:"objectid,attr"`
	Subtype         string            `xml:"subtype,attr"`
	CollID          int64             `xml:"collid,attr"`
	Name            *CollectionName   `xml:"name"`
	OriginalName    string            `xml:"originalname"`
	YearPublished   string            `xml:"yearpublished"`
	Image           string            `xml:"image"`
	Thumbnail       string            `xml:"thumbnail"`
	Stats           *CollectionStats  `xml:"stats"`
	Status          *CollectionStatus `xml:"status"`
	NumPlays        int               `xml:"numplays"`
	Version         *Version          `xml:"version>item"`
	Comment         string            `xml:"comment"`
	WishlistComment string            `xml:"wishlistcomment"`
	ConditionText   string            `xml:"conditiontext"`
}

type CollectionResponse struct {
	XMLName    xml.Name          `xml:"items"`
	TotalItems int               `xml:"totalitems,attr"`
	TermsOfUse string            `xml:"termsofuse,attr"`
	PubDate    string            `xml:"pubdate,attr"`
	Items      []*CollectionItem `xml:"item"`
}

// collectionErrors is returned in place of <items> for an invalid request, e.g. an unknown username
type collectionErrors struct {
	XMLName  xml.Name `xml:"errors"`
	Messages []string `xml:"error>message"`
}
