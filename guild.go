package boardgamegeek

import (
	"encoding/xml"
)

type GuildMemberSort string

const (
	GuildMemberSortUsername GuildMemberSort = "username"
	GuildMemberSortDate     GuildMemberSort = "date"
)

type GuildRequest struct {
	ID      int64 `url:"id"`
	Members bool  `url:"members,int,omitempty"`
	// Sort orders the members list, defaults to username
	Sort GuildMemberSort `url:"sort,omitempty"`
	// Page pages the members list, 25 members per page
	Page int `url:"page,omitempty"`
}

type GuildLocation struct {
	XMLName         xml.Name `xml:"location"`
	Addr1           string   `xml:"addr1"`
	Addr2           string   `xml:"addr2"`
	City            string   `xml:"city"`
	StateOrProvince string   `xml:"stateorprovince"`
	PostalCode      string   `xml:"postalcode"`
	Country         string   `xml:"country"`
}

type GuildMember struct {
	XMLName xml.Name `xml:"member"`
	Name    string   `xml:"name,attr"`
	Date    string   `xml:"date,attr"`
}

type GuildMembers struct {
	XMLName xml.Name       `xml:"members"`
	Count   int            `xml:"count,attr"`
	Page    int            `xml:"page,attr"`
	Members []*GuildMember `xml:"member"`
}

type GuildResponse struct {
	XMLName     xml.Name       `xml:"guild"`
	Id          int64          `xml:"id,attr"`
	Name        string         `xml:"name,attr"`
	Created     string         `xml:"created,attr"`
	TermsOfUse  string         `xml:"termsofuse,attr"`
	Category    string         `xml:"category"`
	Website     string         `xml:"website"`
	Manager     string         `xml:"manager"`
	Description string         `xml:"description"`
	Location    *GuildLocation `xml:"location"`
	Members     *GuildMembers  `xml:"members"`
}
