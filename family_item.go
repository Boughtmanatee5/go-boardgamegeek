package boardgamegeek

import (
	"encoding/xml"
)

// FamilyItemType is the type of a family, used to filter GetFamilyItems
type FamilyItemType string

// Family types accepted by FamilyItemsOptions.Types
const (
	FamilyItemTypeRPG             FamilyItemType = "rpg"
	FamilyItemTypeRPGPeriodical   FamilyItemType = "rpgperiodical"
	FamilyItemTypeBoardGameFamily FamilyItemType = "boardgamefamily"
)

// FamilyItemsOptions are the parameters for GetFamilyItems
type FamilyItemsOptions struct {
	IDs   []int64          `url:"id,comma"`
	Types []FamilyItemType `url:"type,comma,omitempty"`
}

// FamilyLink links a family to a thing in it
type FamilyLink struct {
	XMLName xml.Name `xml:"link"`
	Type    string   `xml:"type,attr"`
	ID      int64    `xml:"id,attr"`
	Value   string   `xml:"value,attr"`
	Inbound bool     `xml:"inbound,attr"`
}

// FamilyItem is a family of related things
type FamilyItem struct {
	XMLName     xml.Name      `xml:"item"`
	Type        string        `xml:"type,attr"`
	ID          int64         `xml:"id,attr"`
	Thumbnail   string        `xml:"thumbnail"`
	Image       string        `xml:"image"`
	Names       []*Name       `xml:"name"`
	Description string        `xml:"description"`
	Links       []*FamilyLink `xml:"link"`
}

// FamilyItemsResponse is the response from GetFamilyItems
type FamilyItemsResponse struct {
	XMLName xml.Name      `xml:"items"`
	Items   []*FamilyItem `xml:"item"`
}
