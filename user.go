package boardgamegeek

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type UserDomain string

const (
	UserDomainBoardGame UserDomain = "boardgame"
	UserDomainRPG       UserDomain = "rpg"
	UserDomainVideoGame UserDomain = "videogame"
)

type UserRequest struct {
	Name    string
	Buddies bool
	Guilds  bool
	Hot     bool
	Top     bool
	// Domain selects which hot and top lists are returned, defaults to boardgame
	Domain UserDomain
	// Page pages the buddies and guilds lists
	Page int
}

func (r *UserRequest) toQueryString() string {
	query := fmt.Sprintf("name=%s", url.QueryEscape(r.Name))

	if r.Buddies {
		query = fmt.Sprintf("%s&buddies=%d", query, 1)
	}

	if r.Guilds {
		query = fmt.Sprintf("%s&guilds=%d", query, 1)
	}

	if r.Hot {
		query = fmt.Sprintf("%s&hot=%d", query, 1)
	}

	if r.Top {
		query = fmt.Sprintf("%s&top=%d", query, 1)
	}

	if r.Domain != "" {
		query = fmt.Sprintf("%s&domain=%s", query, r.Domain)
	}

	if r.Page > 0 {
		query = fmt.Sprintf("%s&page=%d", query, r.Page)
	}

	return query
}

type UserBuddy struct {
	XMLName xml.Name `xml:"buddy"`
	Id      int64    `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

type UserBuddies struct {
	XMLName xml.Name     `xml:"buddies"`
	Total   int          `xml:"total,attr"`
	Page    int          `xml:"page,attr"`
	Buddies []*UserBuddy `xml:"buddy"`
}

type UserGuild struct {
	XMLName xml.Name `xml:"guild"`
	Id      int64    `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

type UserGuilds struct {
	XMLName xml.Name     `xml:"guilds"`
	Total   int          `xml:"total,attr"`
	Page    int          `xml:"page,attr"`
	Guilds  []*UserGuild `xml:"guild"`
}

type UserRankedItem struct {
	XMLName xml.Name `xml:"item"`
	Rank    int      `xml:"rank,attr"`
	Type    string   `xml:"type,attr"`
	Id      int64    `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

// UserRankedItems is a user's hot or top list
type UserRankedItems struct {
	Domain string            `xml:"domain,attr"`
	Items  []*UserRankedItem `xml:"item"`
}

type UserResponse struct {
	XMLName          xml.Name         `xml:"user"`
	Id               int64            `xml:"id,attr"`
	Name             string           `xml:"name,attr"`
	TermsOfUse       string           `xml:"termsofuse,attr"`
	FirstName        ValueAttr        `xml:"firstname"`
	LastName         ValueAttr        `xml:"lastname"`
	AvatarLink       ValueAttr        `xml:"avatarlink"`
	YearRegistered   ValueAttr        `xml:"yearregistered"`
	LastLogin        ValueAttr        `xml:"lastlogin"`
	StateOrProvince  ValueAttr        `xml:"stateorprovince"`
	Country          ValueAttr        `xml:"country"`
	WebAddress       ValueAttr        `xml:"webaddress"`
	XboxAccount      ValueAttr        `xml:"xboxaccount"`
	WiiAccount       ValueAttr        `xml:"wiiaccount"`
	PSNAccount       ValueAttr        `xml:"psnaccount"`
	BattleNetAccount ValueAttr        `xml:"battlenetaccount"`
	SteamAccount     ValueAttr        `xml:"steamaccount"`
	TradeRating      ValueAttr        `xml:"traderating"`
	Buddies          *UserBuddies     `xml:"buddies"`
	Guilds           *UserGuilds      `xml:"guilds"`
	Hot              *UserRankedItems `xml:"hot"`
	Top              *UserRankedItems `xml:"top"`
}
