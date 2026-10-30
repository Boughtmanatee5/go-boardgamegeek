package boardgamegeek

import (
	"encoding/xml"
	"time"
)

// PlaysType is the kind of object PlaysOptions.ID refers to
type PlaysType string

// Object types accepted by PlaysOptions.Type
const (
	PlaysTypeThing  PlaysType = "thing"
	PlaysTypeFamily PlaysType = "family"
)

// PlaysOptions are the parameters for GetPlays and AllPlays. Username or ID is required;
// both can be given to list one user's plays of one thing.
type PlaysOptions struct {
	Username string `url:"username,omitempty"`
	ID       int64  `url:"id,omitempty"`
	// Type is what ID refers to; BGG treats ID as a thing when it is empty
	Type PlaysType `url:"type,omitempty"`
	// MinDate and MaxDate limit plays to a date range, inclusive; zero means no limit
	MinDate time.Time `url:"mindate,omitempty" layout:"2006-01-02"`
	MaxDate time.Time `url:"maxdate,omitempty" layout:"2006-01-02"`
	Subtype ThingType `url:"subtype,omitempty"`
	// Page selects a page of 100 plays, starting at 1
	Page int `url:"page,omitempty"`
}

// Player is one player recorded in a play. Players without a BGG account have an empty Username and UserID 0.
type Player struct {
	XMLName  xml.Name `xml:"player"`
	Username string   `xml:"username,attr"`
	UserID   int64    `xml:"userid,attr"`
	Name     string   `xml:"name,attr"`
	// StartPosition, Color and Score are free text as entered by the logger, and often empty
	StartPosition string `xml:"startposition,attr"`
	Color         string `xml:"color,attr"`
	Score         string `xml:"score,attr"`
	New           bool   `xml:"new,attr"`
	Rating        int    `xml:"rating,attr"`
	Win           bool   `xml:"win,attr"`
}

// PlayItem is the thing that was played
type PlayItem struct {
	XMLName    xml.Name `xml:"item"`
	Name       string   `xml:"name,attr"`
	ObjectType string   `xml:"objecttype,attr"`
	ObjectID   int64    `xml:"objectid,attr"`
	// Subtypes are the thing's types, e.g. boardgame and boardgameimplementation
	Subtypes []ValueAttr `xml:"subtypes>subtype"`
}

// Play is one logged play
type Play struct {
	XMLName xml.Name `xml:"play"`
	ID      int64    `xml:"id,attr"`
	// UserID is the user who logged the play; BGG only sends it when plays are listed by ID
	UserID int64 `xml:"userid,attr"`
	// Date is the play date as entered, YYYY-MM-DD; BGG does not validate it
	Date     string `xml:"date,attr"`
	Quantity int    `xml:"quantity,attr"`
	// Length is the play time in minutes, 0 if not recorded
	Length     int       `xml:"length,attr"`
	Incomplete bool      `xml:"incomplete,attr"`
	NowInStats bool      `xml:"nowinstats,attr"`
	Location   string    `xml:"location,attr"`
	Item       *PlayItem `xml:"item"`
	Comments   string    `xml:"comments"`
	Players    []*Player `xml:"players>player"`
}

// PlaysResponse is the response from GetPlays
type PlaysResponse struct {
	XMLName xml.Name `xml:"plays"`
	// Username and UserID are set when plays are listed by Username
	Username string  `xml:"username,attr"`
	UserID   int64   `xml:"userid,attr"`
	Total    int     `xml:"total,attr"`
	Page     int     `xml:"page,attr"`
	Plays    []*Play `xml:"play"`
}

// messageBox is the html fragment, sent after an xml declaration with HTTP 200, that /plays returns
// for an unknown user or object: <div class='messagebox error'>Invalid object or user</div>
type messageBox struct {
	XMLName xml.Name `xml:"div"`
	Class   string   `xml:"class,attr"`
	Message string   `xml:",chardata"`
}
