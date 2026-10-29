package boardgamegeek

import (
	"encoding/xml"
)

// ThingType is the type of a thing, used to filter Search and GetThings
type ThingType string

// Thing types accepted by SearchOptions.Types and ThingsOptions.Types
const (
	ThingTypeRPGItem            ThingType = "rpgitem"
	ThingTypeVideoGame          ThingType = "videogame"
	ThingTypeBoardGame          ThingType = "boardgame"
	ThingTypeBoardGameAccessory ThingType = "boardgameaccessory"
	ThingTypeBoardGameExpansion ThingType = "boardgameexpansion"
	ThingTypeBoardGameDesigner  ThingType = "boardgamedesigner"
)

// ThingsOptions are the parameters for GetThings and AllThingComments
type ThingsOptions struct {
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

// Statistics holds a thing's rating statistics, returned with ThingsOptions.Stats
type Statistics struct {
	XMLName xml.Name `xml:"statistics"`
	Page    int      `xml:"page,attr"`
}

// Version is one published version of a thing, returned with ThingsOptions.Versions
type Version struct {
	XMLName xml.Name `xml:"item"`
	Type    string   `xml:"type,attr"`
	ID      int64    `xml:"id,attr"`
}

// Versions lists a thing's published versions
type Versions struct {
	XMLName xml.Name   `xml:"versions"`
	Items   []*Version `xml:"item"`
}

// Video is one video linked to a thing, returned with ThingsOptions.Videos
type Video struct {
	XMLName xml.Name `xml:"video"`
	ID      int64    `xml:"id,attr"`
}

// Videos lists a thing's videos
type Videos struct {
	XMLName xml.Name `xml:"videos"`
	Total   int64    `xml:"total,attr"`
	Videos  []*Video `xml:"video"`
}

// Comment is one user comment or rating on a thing
type Comment struct {
	XMLName  xml.Name `xml:"comment"`
	Username string   `xml:"username,attr"`
	Rating   string   `xml:"rating,attr"`
	Value    string   `xml:"value,attr"`
}

// Comments is one page of a thing's comments, returned with ThingsOptions.Comments or RatingComments
type Comments struct {
	XMLName    xml.Name   `xml:"comments"`
	Page       int        `xml:"page,attr"`
	TotalItems int64      `xml:"totalitems,attr"`
	Comments   []*Comment `xml:"comment"`
}

// MarketplaceListing is one GeekMarket listing for a thing, returned with ThingsOptions.MarketPlace
type MarketplaceListing struct {
	XMLName xml.Name `xml:"listing"`
}

// MarketplaceListings lists a thing's GeekMarket listings
type MarketplaceListings struct {
	XMLName  xml.Name              `xml:"marketplacelistings"`
	Listings []*MarketplaceListing `xml:"listing"`
}

// Thing is a board game, expansion, accessory, RPG item or video game
type Thing struct {
	XMLName             xml.Name             `xml:"item"`
	Type                string               `xml:"type,attr"`
	ID                  int64                `xml:"id,attr"`
	Names               []*Name              `xml:"name"`
	YearPublished       *ValueAttr           `xml:"yearpublished"`
	Statistics          *Statistics          `xml:"statistics"`
	Versions            *Versions            `xml:"versions"`
	Videos              *Videos              `xml:"videos"`
	Comments            *Comments            `xml:"comments"`
	MarketplaceListings *MarketplaceListings `xml:"marketplacelistings"`
}

// ThingsResponse is the response from GetThings
type ThingsResponse struct {
	XMLName xml.Name `xml:"items"`
	Items   []*Thing `xml:"item"`
}
