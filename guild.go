package boardgamegeek

import (
	"encoding/xml"
)

// GuildMemberSort orders the members GetGuild returns
type GuildMemberSort string

// Sort orders accepted by GuildOptions.Sort
const (
	GuildMemberSortUsername GuildMemberSort = "username"
	GuildMemberSortDate     GuildMemberSort = "date"
)

// GuildOptions are the parameters for GetGuild and AllGuildMembers
type GuildOptions struct {
	ID      int64 `url:"id"`
	Members bool  `url:"members,int,omitempty"`
	// Sort orders the members list, defaults to username
	Sort GuildMemberSort `url:"sort,omitempty"`
	// Page pages the members list, 25 members per page
	Page int `url:"page,omitempty"`
}

// GuildLocation is a guild's address
type GuildLocation struct {
	XMLName         xml.Name `xml:"location"`
	Addr1           string   `xml:"addr1"`
	Addr2           string   `xml:"addr2"`
	City            string   `xml:"city"`
	StateOrProvince string   `xml:"stateorprovince"`
	PostalCode      string   `xml:"postalcode"`
	Country         string   `xml:"country"`
}

// GuildMember is one member of a guild
type GuildMember struct {
	XMLName xml.Name `xml:"member"`
	Name    string   `xml:"name,attr"`
	Date    string   `xml:"date,attr"`
}

// GuildMembers is one page of a guild's members, returned with GuildOptions.Members
type GuildMembers struct {
	XMLName xml.Name       `xml:"members"`
	Count   int            `xml:"count,attr"`
	Page    int            `xml:"page,attr"`
	Members []*GuildMember `xml:"member"`
}

// GuildResponse is the response from GetGuild
type GuildResponse struct {
	XMLName     xml.Name       `xml:"guild"`
	ID          int64          `xml:"id,attr"`
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
