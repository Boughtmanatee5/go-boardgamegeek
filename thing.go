package boardgamegeek

import (
	"encoding/xml"
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
	IDs            []int64     `url:"id,comma"`
	Types          []ThingType `url:"type,comma,omitempty"`
	Versions       bool        `url:"versions,int,omitempty"`
	Videos         bool        `url:"videos,int,omitempty"`
	Stats          bool        `url:"stats,int,omitempty"`
	MarketPlace    bool        `url:"marketplace,int,omitempty"`
	Comments       bool        `url:"comments,int,omitempty"`
	RatingComments bool        `url:"ratingcomments,int,omitempty"`
	Page           int         `url:"page,omitempty"`
	PageSize       int         `url:"pagesize,omitempty"`
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
	YearPublished       *ValueAttr           `xml:"yearpublished"`
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
