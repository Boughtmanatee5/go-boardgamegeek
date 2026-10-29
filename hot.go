package boardgamegeek

import (
	"encoding/xml"
)

// HotItemType selects which hotness list GetHotItems returns
type HotItemType string

// Hotness lists accepted by HotItemsOptions.Type
const (
	HotItemTypeBoardGame        HotItemType = "boardgame"
	HotItemTypeRPG              HotItemType = "rpg"
	HotItemTypeVideoGame        HotItemType = "videogame"
	HotItemTypeBoardGamePerson  HotItemType = "boardgameperson"
	HotItemTypeRPGPerson        HotItemType = "rpgperson"
	HotItemTypeBoardGameCompany HotItemType = "boardgamecompany"
	HotItemTypeRPGCompany       HotItemType = "rpgcompany"
	HotItemTypeVideoGameCompany HotItemType = "videogamecompany"
)

// HotItemsOptions takes a single Type: BGG ignores a list of types and falls back to boardgame.
type HotItemsOptions struct {
	Type HotItemType `url:"type,omitempty"`
}

// HotItem is one entry in a hotness list
type HotItem struct {
	XMLName   xml.Name  `xml:"item"`
	ID        int64     `xml:"id,attr"`
	Rank      int       `xml:"rank,attr"`
	Thumbnail ValueAttr `xml:"thumbnail"`
	Name      ValueAttr `xml:"name"`
	// Only present for game types, not people or companies
	YearPublished *ValueAttr `xml:"yearpublished"`
}

// HotItemsResponse is the response from GetHotItems
type HotItemsResponse struct {
	XMLName    xml.Name   `xml:"items"`
	TermsOfUse string     `xml:"termsofuse,attr"`
	Items      []*HotItem `xml:"item"`
}
