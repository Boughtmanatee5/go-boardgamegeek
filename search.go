package boardgamegeek

import (
	"encoding/xml"
)

// SearchType is the type of item to search for: every ThingType, plus types such as
// boardgamedesigner that /search accepts but /thing does not
type SearchType string

// Search types accepted by SearchOptions.Types
const (
	SearchTypeRPGItem            SearchType = "rpgitem"
	SearchTypeRPGIssue           SearchType = "rpgissue"
	SearchTypeVideoGame          SearchType = "videogame"
	SearchTypeBoardGame          SearchType = "boardgame"
	SearchTypeBoardGameAccessory SearchType = "boardgameaccessory"
	SearchTypeBoardGameExpansion SearchType = "boardgameexpansion"
	SearchTypeBoardGameDesigner  SearchType = "boardgamedesigner"
)

// SearchOptions are the parameters for Search
type SearchOptions struct {
	Query string       `url:"query"`
	Types []SearchType `url:"type,comma,omitempty"`
	Exact bool         `url:"exact,int,omitempty"`
}

// ValueAttr unmarshals an element whose only content is a value attribute, e.g. <country value="United States" />,
// to the attribute's value. Use *ValueAttr for optional elements so they are nil when absent.
type ValueAttr string

// UnmarshalXML implements xml.Unmarshaler by reading the element's value attribute
func (v *ValueAttr) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "value" {
			*v = ValueAttr(attr.Value)
		}
	}
	return d.Skip()
}

// Name is one of an item's names; Type is "primary" or "alternate"
type Name struct {
	XMLName   xml.Name `xml:"name"`
	Type      string   `xml:"type,attr"`
	SortIndex int      `xml:"sortindex,attr"`
	Value     string   `xml:"value,attr"`
}

// SearchItem is one match in a SearchResponse
type SearchItem struct {
	XMLName       xml.Name   `xml:"item"`
	Text          string     `xml:",chardata"`
	Type          string     `xml:"type,attr"`
	ID            int64      `xml:"id,attr"`
	Name          *Name      `xml:"name"`
	YearPublished *ValueAttr `xml:"yearpublished"`
}

// SearchResponse is the response from Search
type SearchResponse struct {
	XMLName xml.Name      `xml:"items"`
	Total   int64         `xml:"total,attr"`
	Items   []*SearchItem `xml:"item"`
}
