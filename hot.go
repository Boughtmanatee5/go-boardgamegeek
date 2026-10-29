package boardgamegeek

import (
	"encoding/xml"
	"fmt"
)

type HotItemType string

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

// HotItemsRequest takes a single Type: BGG ignores a list of types and falls back to boardgame.
type HotItemsRequest struct {
	Type HotItemType
}

func (r *HotItemsRequest) toQueryString() string {
	if r.Type == "" {
		return ""
	}

	return fmt.Sprintf("type=%s", r.Type)
}

type HotItem struct {
	XMLName   xml.Name  `xml:"item"`
	Id        int64     `xml:"id,attr"`
	Rank      int       `xml:"rank,attr"`
	Thumbnail ValueAttr `xml:"thumbnail"`
	Name      ValueAttr `xml:"name"`
	// Only present for game types, not people or companies
	YearPublished *ValueAttr `xml:"yearpublished"`
}

type HotItemsResponse struct {
	XMLName    xml.Name   `xml:"items"`
	TermsOfUse string     `xml:"termsofuse,attr"`
	Items      []*HotItem `xml:"item"`
}
