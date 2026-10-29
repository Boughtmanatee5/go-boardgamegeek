package boardgamegeek

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type FamilyItemType string

const (
	FamilyItemTypeRPG             FamilyItemType = "rpg"
	FamilyItemTypeRPGPeriodical   FamilyItemType = "rpgperiodical"
	FamilyItemTypeBoardGameFamily FamilyItemType = "boardgamefamily"
)

type FamilyItemRequest struct {
	IDs   []int64
	Types []FamilyItemType
}

func (r *FamilyItemRequest) toQueryString() string {
	var ids []string
	for _, id := range r.IDs {
		ids = append(ids, strconv.FormatInt(id, 10))
	}
	query := fmt.Sprintf("id=%s", strings.Join(ids, ","))

	if len(r.Types) != 0 {
		var types []string
		for _, familyType := range r.Types {
			types = append(types, string(familyType))
		}
		query = fmt.Sprintf("%s&type=%s", query, strings.Join(types, ","))
	}

	return query
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
