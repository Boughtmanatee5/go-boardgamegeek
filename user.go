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
	FirstName        ValueElement     `xml:"firstname"`
	LastName         ValueElement     `xml:"lastname"`
	AvatarLink       ValueElement     `xml:"avatarlink"`
	YearRegistered   ValueElement     `xml:"yearregistered"`
	LastLogin        ValueElement     `xml:"lastlogin"`
	StateOrProvince  ValueElement     `xml:"stateorprovince"`
	Country          ValueElement     `xml:"country"`
	WebAddress       ValueElement     `xml:"webaddress"`
	XboxAccount      ValueElement     `xml:"xboxaccount"`
	WiiAccount       ValueElement     `xml:"wiiaccount"`
	PSNAccount       ValueElement     `xml:"psnaccount"`
	BattleNetAccount ValueElement     `xml:"battlenetaccount"`
	SteamAccount     ValueElement     `xml:"steamaccount"`
	TradeRating      ValueElement     `xml:"traderating"`
	Buddies          *UserBuddies     `xml:"buddies"`
	Guilds           *UserGuilds      `xml:"guilds"`
	Hot              *UserRankedItems `xml:"hot"`
	Top              *UserRankedItems `xml:"top"`
}
