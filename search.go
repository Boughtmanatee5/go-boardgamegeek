package boardgamegeek

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type SearchRequest struct {
	Query string
	Types []ThingType
	Exact bool
}

func (sr *SearchRequest) toQueryString() string {
	searchQuery := fmt.Sprintf("query=%s", strings.ReplaceAll(sr.Query, " ", "+"))

	if len(sr.Types) != 0 {
		typesQuery := "type="
		for i, bggType := range sr.Types {
			if i == 0 {
				typesQuery = fmt.Sprintf("%s%s", typesQuery, bggType)
				continue
			}
			typesQuery = fmt.Sprintf("%s,%s", typesQuery, bggType)
		}
		searchQuery = fmt.Sprintf("%s&%s", searchQuery, typesQuery)
	}

	if sr.Exact {
		searchQuery = fmt.Sprintf("%s&exact=%d", searchQuery, 1)
	}

	return searchQuery
}

// ValueAttr unmarshals an element whose only content is a value attribute, e.g. <country value="United States" />,
// to the attribute's value. Use *ValueAttr for optional elements so they are nil when absent.
type ValueAttr string

func (v *ValueAttr) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "value" {
			*v = ValueAttr(attr.Value)
		}
	}
	return d.Skip()
}

type Name struct {
	XMLName   xml.Name `xml:"name"`
	Type      string   `xml:"type,attr"`
	SortIndex int      `xml:"sortindex,attr"`
	Value     string   `xml:"value,attr"`
}

type SearchItem struct {
	XMLName       xml.Name   `xml:"item"`
	Text          string     `xml:",chardata"`
	Type          string     `xml:"type,attr"`
	Id            int64      `xml:"id,attr"`
	Name          *Name      `xml:"name"`
	YearPublished *ValueAttr `xml:"yearpublished"`
}

type SearchResponse struct {
	XMLName xml.Name      `xml:"items"`
	Total   int64         `xml:"total,attr"`
	Items   []*SearchItem `xml:"item"`
}
