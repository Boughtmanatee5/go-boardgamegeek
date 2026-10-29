package boardgamegeek

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ThingType string

const (
	RPGITEM            ThingType = "rpgitem"
	VIDEOGAME          ThingType = "videogame"
	BOARDGAME          ThingType = "boardgame"
	BOARDGAMEACCESSORY ThingType = "boardgameaccessory"
	BOARDGAMEEXPANSION ThingType = "boardgameexpansion"
	BOARDGAMEDESIGNER  ThingType = "boardgamedesigner"
)

type GetThingsRequest struct {
	IDs            []int64
	Types          []ThingType
	Versions       bool
	Videos         bool
	Stats          bool
	MarketPlace    bool
	Comments       bool
	RatingComments bool
	Page           int
	PageSize       int
}

func (r *GetThingsRequest) toQueryString() string {
	var ids []string
	for _, id := range r.IDs {
		ids = append(ids, strconv.FormatInt(id, 10))
	}
	query := fmt.Sprintf("id=%s", strings.Join(ids, ","))

	if len(r.Types) != 0 {
		typesQuery := "type="
		for i, bggType := range r.Types {
			if i == 0 {
				typesQuery = fmt.Sprintf("%s%s", typesQuery, bggType)
				continue
			}
			typesQuery = fmt.Sprintf("%s,%s", typesQuery, bggType)
		}
		query = fmt.Sprintf("%s&%s", query, typesQuery)
	}

	if r.Versions {
		query = fmt.Sprintf("%s&versions=%d", query, 1)
	}

	if r.Videos {
		query = fmt.Sprintf("%s&videos=%d", query, 1)
	}

	if r.Stats {
		query = fmt.Sprintf("%s&stats=%d", query, 1)
	}

	if r.MarketPlace {
		query = fmt.Sprintf("%s&marketplace=%d", query, 1)
	}

	if r.Comments {
		query = fmt.Sprintf("%s&comments=%d", query, 1)
	}

	if r.RatingComments {
		query = fmt.Sprintf("%s&ratingcomments=%d", query, 1)
	}

	if r.Page > 0 {
		query = fmt.Sprintf("%s&page=%d", query, r.Page)
	}

	if r.PageSize > 0 {
		query = fmt.Sprintf("%s&pagesize=%d", query, r.PageSize)
	}

	return query
}

type Statistics struct {
	XMLName xml.Name `xml:"statistics"`
	Page    int      `xml:"page,attr"`
}

type Version struct {
	XMLName xml.Name `xml:"item"`
	Type    string   `xml:"type,attr"`
	Id      int64    `xml:"id,attr"`
}

type Versions struct {
	XMLName xml.Name   `xml:"versions"`
	Items   []*Version `xml:"item"`
}

type Video struct {
	XMLName xml.Name `xml:"video"`
	Id      int64    `xml:"id,attr"`
}

type Videos struct {
	XMLName xml.Name `xml:"videos"`
	Total   int64    `xml:"total,attr"`
	Videos  []*Video `xml:"video"`
}

type Comment struct {
	XMLName  xml.Name `xml:"comment"`
	Username string   `xml:"username,attr"`
	Rating   string   `xml:"rating,attr"`
	Value    string   `xml:"value,attr"`
}

type Comments struct {
	XMLName    xml.Name   `xml:"comments"`
	Page       int        `xml:"page,attr"`
	TotalItems int64      `xml:"totalitems,attr"`
	Comments   []*Comment `xml:"comment"`
}

type MarketplaceListing struct {
	XMLName xml.Name `xml:"listing"`
}

type MarketplaceListings struct {
	XMLName  xml.Name              `xml:"marketplacelistings"`
	Listings []*MarketplaceListing `xml:"listing"`
}

type Thing struct {
	XMLName             xml.Name             `xml:"item"`
	Type                string               `xml:"type,attr"`
	Id                  int64                `xml:"id,attr"`
	Names               []*Name              `xml:"name"`
	YearPublished       *YearPublished       `xml:"yearpublished"`
	Statistics          *Statistics          `xml:"statistics"`
	Versions            *Versions            `xml:"versions"`
	Videos              *Videos              `xml:"videos"`
	Comments            *Comments            `xml:"comments"`
	MarketplaceListings *MarketplaceListings `xml:"marketplacelistings"`
}

type GetThingsResponse struct {
	XMLName xml.Name `xml:"items"`
	Items   []*Thing `xml:"item"`
}
