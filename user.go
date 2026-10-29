package boardgamegeek

import (
	"encoding/xml"
)

type UserDomain string

const (
	UserDomainBoardGame UserDomain = "boardgame"
	UserDomainRPG       UserDomain = "rpg"
	UserDomainVideoGame UserDomain = "videogame"
)

type UserOptions struct {
	Name    string `url:"name"`
	Buddies bool   `url:"buddies,int,omitempty"`
	Guilds  bool   `url:"guilds,int,omitempty"`
	Hot     bool   `url:"hot,int,omitempty"`
	Top     bool   `url:"top,int,omitempty"`
	// Domain selects which hot and top lists are returned, defaults to boardgame
	Domain UserDomain `url:"domain,omitempty"`
	// Page pages the buddies and guilds lists
	Page int `url:"page,omitempty"`
}

type UserBuddy struct {
	XMLName xml.Name `xml:"buddy"`
	ID      int64    `xml:"id,attr"`
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
	ID      int64    `xml:"id,attr"`
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
	ID      int64    `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

// UserRankedItems is a user's hot or top list
type UserRankedItems struct {
	Domain string            `xml:"domain,attr"`
	Items  []*UserRankedItem `xml:"item"`
}

type UserResponse struct {
	XMLName          xml.Name         `xml:"user"`
	ID               int64            `xml:"id,attr"`
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
