package boardgamegeek

import (
	"encoding/xml"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

// ErrCollectionNotReady matches, with errors.Is, the *AcceptedError returned when BGG has queued
// the collection request (HTTP 202). Retry the same request after a few seconds.
var ErrCollectionNotReady = errors.New("collection request queued by boardgamegeek, retry later")

// Bool returns a pointer to b, for the optional filters on CollectionRequest.
func Bool(b bool) *bool {
	return &b
}

type CollectionRequest struct {
	// Username is required
	Username       string
	Version        bool
	Subtype        ThingType
	ExcludeSubtype ThingType
	IDs            []int64
	// Brief returns only the name and status of each item
	Brief bool
	Stats bool
	// Status filters: nil means no filter, true includes only matching items, false excludes them
	Own        *bool
	Rated      *bool
	Played     *bool
	Comment    *bool
	Trade      *bool
	Want       *bool
	Wishlist   *bool
	Preordered *bool
	WantToPlay *bool
	WantToBuy  *bool
	PrevOwned  *bool
	HasParts   *bool
	WantParts  *bool
	// WishlistPriority filters on wishlist priority 1-5, 0 means no filter
	WishlistPriority int
	// Rating filters are 1-10, 0 means no filter
	MinRating    float64
	Rating       float64
	MinBGGRating float64
	BGGRating    float64
	// Play count filters, 0 means no filter
	MinPlays int
	MaxPlays int
	// ShowPrivate only returns private info for the authenticated user's own collection
	ShowPrivate bool
	CollID      int64
	// ModifiedSince is YY-MM-DD or YY-MM-DD HH:MM:SS
	ModifiedSince string
}

func (r *CollectionRequest) toQueryString() string {
	query := url.Values{}
	query.Set("username", r.Username)

	if r.Version {
		query.Set("version", "1")
	}
	if r.Subtype != "" {
		query.Set("subtype", string(r.Subtype))
	}
	if r.ExcludeSubtype != "" {
		query.Set("excludesubtype", string(r.ExcludeSubtype))
	}
	if len(r.IDs) != 0 {
		var ids []string
		for _, id := range r.IDs {
			ids = append(ids, strconv.FormatInt(id, 10))
		}
		query.Set("id", strings.Join(ids, ","))
	}
	if r.Brief {
		query.Set("brief", "1")
	}
	if r.Stats {
		query.Set("stats", "1")
	}

	filters := map[string]*bool{
		"own":        r.Own,
		"rated":      r.Rated,
		"played":     r.Played,
		"comment":    r.Comment,
		"trade":      r.Trade,
		"want":       r.Want,
		"wishlist":   r.Wishlist,
		"preordered": r.Preordered,
		"wanttoplay": r.WantToPlay,
		"wanttobuy":  r.WantToBuy,
		"prevowned":  r.PrevOwned,
		"hasparts":   r.HasParts,
		"wantparts":  r.WantParts,
	}
	for param, filter := range filters {
		if filter == nil {
			continue
		}
		if *filter {
			query.Set(param, "1")
		} else {
			query.Set(param, "0")
		}
	}

	if r.WishlistPriority > 0 {
		query.Set("wishlistpriority", strconv.Itoa(r.WishlistPriority))
	}
	if r.MinRating > 0 {
		query.Set("minrating", strconv.FormatFloat(r.MinRating, 'f', -1, 64))
	}
	if r.Rating > 0 {
		query.Set("rating", strconv.FormatFloat(r.Rating, 'f', -1, 64))
	}
	if r.MinBGGRating > 0 {
		query.Set("minbggrating", strconv.FormatFloat(r.MinBGGRating, 'f', -1, 64))
	}
	if r.BGGRating > 0 {
		query.Set("bggrating", strconv.FormatFloat(r.BGGRating, 'f', -1, 64))
	}
	if r.MinPlays > 0 {
		query.Set("minplays", strconv.Itoa(r.MinPlays))
	}
	if r.MaxPlays > 0 {
		query.Set("maxplays", strconv.Itoa(r.MaxPlays))
	}
	if r.ShowPrivate {
		query.Set("showprivate", "1")
	}
	if r.CollID > 0 {
		query.Set("collid", strconv.FormatInt(r.CollID, 10))
	}
	if r.ModifiedSince != "" {
		query.Set("modifiedsince", r.ModifiedSince)
	}

	return query.Encode()
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
	Id           int64    `xml:"id,attr"`
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
	ObjectId        int64             `xml:"objectid,attr"`
	Subtype         string            `xml:"subtype,attr"`
	CollId          int64             `xml:"collid,attr"`
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
