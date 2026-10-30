package boardgamegeek

import (
	"encoding/xml"
	"fmt"
	"strconv"
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
	*v = ValueAttr(valueAttr(start))
	return d.Skip()
}

// valueAttr returns the value attribute of start
func valueAttr(start xml.StartElement) string {
	for _, attr := range start.Attr {
		if attr.Name.Local == "value" {
			return attr.Value
		}
	}
	return ""
}

// IntValue unmarshals an element whose only content is a numeric value attribute, e.g. <minplayers value="3" />,
// to the attribute's value. An empty value is 0. Use *IntValue for optional elements so they are nil when absent.
type IntValue int

// UnmarshalXML implements xml.Unmarshaler by parsing the element's value attribute as an integer
func (v *IntValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if value := valueAttr(start); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("<%s value=%q>: %w", start.Name.Local, value, err)
		}
		*v = IntValue(n)
	}
	return d.Skip()
}

// FloatValue unmarshals an element whose only content is a decimal value attribute, e.g. <average value="7.1" />,
// to the attribute's value. An empty value is 0. Use *FloatValue for optional elements so they are nil when absent.
type FloatValue float64

// UnmarshalXML implements xml.Unmarshaler by parsing the element's value attribute as a float
func (v *FloatValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if value := valueAttr(start); value != "" {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("<%s value=%q>: %w", start.Name.Local, value, err)
		}
		*v = FloatValue(f)
	}
	return d.Skip()
}

// Link links an item to a related object, e.g. a thing to its designer, category or expansion.
// Inbound is true when the link points from the related object to this item, e.g. from an expansion to its base game.
type Link struct {
	XMLName xml.Name `xml:"link"`
	Type    string   `xml:"type,attr"`
	ID      int64    `xml:"id,attr"`
	Value   string   `xml:"value,attr"`
	Inbound bool     `xml:"inbound,attr"`
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
	XMLName       xml.Name  `xml:"item"`
	Text          string    `xml:",chardata"`
	Type          string    `xml:"type,attr"`
	ID            int64     `xml:"id,attr"`
	Name          *Name     `xml:"name"`
	YearPublished *IntValue `xml:"yearpublished"`
}

// SearchResponse is the response from Search
type SearchResponse struct {
	XMLName xml.Name      `xml:"items"`
	Total   int64         `xml:"total,attr"`
	Items   []*SearchItem `xml:"item"`
}
