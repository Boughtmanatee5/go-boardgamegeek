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

type YearPublished struct {
	XMLName xml.Name `xml:"yearpublished"`
	Value   string   `xml:"value,attr"`
}

// ValueElement is an element whose only content is a value attribute, e.g. <country value="United States" />
type ValueElement struct {
	Value string `xml:"value,attr"`
}

type Name struct {
	XMLName   xml.Name `xml:"name"`
	Type      string   `xml:"type,attr"`
	SortIndex int      `xml:"sortindex,attr"`
	Value     string   `xml:"value,attr"`
}

type SearchItem struct {
	XMLName       xml.Name       `xml:"item"`
	Text          string         `xml:",chardata"`
	Type          string         `xml:"type,attr"`
	Id            int64          `xml:"id,attr"`
	Name          *Name          `xml:"name"`
	YearPublished *YearPublished `xml:"yearpublished"`
}

type SearchResponse struct {
	XMLName xml.Name      `xml:"items"`
	Total   int64         `xml:"total,attr"`
	Items   []*SearchItem `xml:"item"`
}
