package boardgamegeek

import (
	"encoding/xml"
)

type FamilyItemType string

const (
	FamilyItemTypeRPG             FamilyItemType = "rpg"
	FamilyItemTypeRPGPeriodical   FamilyItemType = "rpgperiodical"
	FamilyItemTypeBoardGameFamily FamilyItemType = "boardgamefamily"
)

type FamilyItemRequest struct {
	IDs   []int64          `url:"id,comma"`
	Types []FamilyItemType `url:"type,comma,omitempty"`
}

type FamilyLink struct {
	XMLName xml.Name `xml:"link"`
	Type    string   `xml:"type,attr"`
	Id      int64    `xml:"id,attr"`
	Value   string   `xml:"value,attr"`
	Inbound bool     `xml:"inbound,attr"`
}

type FamilyItem struct {
	XMLName     xml.Name      `xml:"item"`
	Type        string        `xml:"type,attr"`
	Id          int64         `xml:"id,attr"`
	Thumbnail   string        `xml:"thumbnail"`
	Image       string        `xml:"image"`
	Names       []*Name       `xml:"name"`
	Description string        `xml:"description"`
	Links       []*FamilyLink `xml:"link"`
}

type FamilyItemsResponse struct {
	XMLName xml.Name      `xml:"items"`
	Items   []*FamilyItem `xml:"item"`
}
