package boardgamegeek

import (
	"encoding/xml"
	"time"
)

// ThingType is the type of a thing, used to filter GetThings and GetCollection.
// Search takes a SearchType, which also covers types that are not things.
type ThingType string

// Thing types accepted by ThingsOptions.Types and CollectionOptions.Subtype
const (
	ThingTypeRPGItem            ThingType = "rpgitem"
	ThingTypeRPGIssue           ThingType = "rpgissue"
	ThingTypeVideoGame          ThingType = "videogame"
	ThingTypeBoardGame          ThingType = "boardgame"
	ThingTypeBoardGameAccessory ThingType = "boardgameaccessory"
	ThingTypeBoardGameExpansion ThingType = "boardgameexpansion"
)

// ThingsOptions are the parameters for GetThings and AllThingComments
type ThingsOptions struct {
	IDs            []int64     `url:"id,comma"`
	Types          []ThingType `url:"type,comma,omitempty"`
	Versions       bool        `url:"versions,int,omitempty"`
	Videos         bool        `url:"videos,int,omitempty"`
	Stats          bool        `url:"stats,int,omitempty"`
	MarketPlace    bool        `url:"marketplace,int,omitempty"`
	Comments       bool        `url:"comments,int,omitempty"`
	RatingComments bool        `url:"ratingcomments,int,omitempty"`
	Page           int         `url:"page,omitempty"`
	PageSize       int         `url:"pagesize,omitempty"`
}

// PollResult is one answer in a poll and how many users voted for it
type PollResult struct {
	XMLName xml.Name `xml:"result"`
	// Level ranks the answers of the language_dependence poll, 1 (no text) to 5 (unplayable in another language)
	Level    int    `xml:"level,attr"`
	Value    string `xml:"value,attr"`
	NumVotes int    `xml:"numvotes,attr"`
}

// PollResults are the answers to a poll. In the suggested_numplayers poll there is one PollResults
// per player count, such as "3" or "4+", each answered Best, Recommended or Not Recommended.
type PollResults struct {
	XMLName    xml.Name      `xml:"results"`
	NumPlayers string        `xml:"numplayers,attr"`
	Results    []*PollResult `xml:"result"`
}

// Poll is a community poll on a board game: suggested_numplayers, suggested_playerage or language_dependence
type Poll struct {
	XMLName    xml.Name       `xml:"poll"`
	Name       string         `xml:"name,attr"`
	Title      string         `xml:"title,attr"`
	TotalVotes int            `xml:"totalvotes,attr"`
	Results    []*PollResults `xml:"results"`
}

// PollSummaryResult is one line of a PollSummary, e.g. Name "bestwith" and Value "Best with 4 players"
type PollSummaryResult struct {
	XMLName xml.Name `xml:"result"`
	// Name is as BGG sends it, including its spelling "recommmendedwith"
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// PollSummary summarises the suggested_numplayers poll
type PollSummary struct {
	XMLName xml.Name             `xml:"poll-summary"`
	Name    string               `xml:"name,attr"`
	Title   string               `xml:"title,attr"`
	Results []*PollSummaryResult `xml:"result"`
}

// Rank is a thing's position in one BGG ranking, overall (Type "subtype") or within a family such as strategy games
type Rank struct {
	XMLName      xml.Name `xml:"rank"`
	Type         string   `xml:"type,attr"`
	ID           int64    `xml:"id,attr"`
	Name         string   `xml:"name,attr"`
	FriendlyName string   `xml:"friendlyname,attr"`
	// Value and BayesAverage are numbers, "Not Ranked", or empty
	Value        string `xml:"value,attr"`
	BayesAverage string `xml:"bayesaverage,attr"`
}

// Ratings summarises how users have rated, owned and weighted a thing
type Ratings struct {
	XMLName      xml.Name    `xml:"ratings"`
	UsersRated   *IntValue   `xml:"usersrated"`
	Average      *FloatValue `xml:"average"`
	BayesAverage *FloatValue `xml:"bayesaverage"`
	Ranks        []*Rank     `xml:"ranks>rank"`
	StdDev       *FloatValue `xml:"stddev"`
	Median       *FloatValue `xml:"median"`
	Owned        *IntValue   `xml:"owned"`
	Trading      *IntValue   `xml:"trading"`
	Wanting      *IntValue   `xml:"wanting"`
	Wishing      *IntValue   `xml:"wishing"`
	NumComments  *IntValue   `xml:"numcomments"`
	NumWeights   *IntValue   `xml:"numweights"`
	// AverageWeight is the average complexity vote, 1 (light) to 5 (heavy), or 0 with no votes
	AverageWeight *FloatValue `xml:"averageweight"`
}

// Statistics holds a thing's rating statistics, returned with ThingsOptions.Stats
type Statistics struct {
	XMLName xml.Name `xml:"statistics"`
	Page    int      `xml:"page,attr"`
	Ratings *Ratings `xml:"ratings"`
}

// Version is one published version of a thing, returned with ThingsOptions.Versions.
// Fields that do not apply to the thing's type, such as ISBNs on a board game, are nil.
type Version struct {
	XMLName       xml.Name   `xml:"item"`
	Type          string     `xml:"type,attr"`
	ID            int64      `xml:"id,attr"`
	Thumbnail     string     `xml:"thumbnail"`
	Image         string     `xml:"image"`
	Names         []*Name    `xml:"name"`
	CanonicalName *ValueAttr `xml:"canonicalname"`
	// Links include the version's publishers, artists and languages, and an inbound link to the thing
	Links         []*Link    `xml:"link"`
	YearPublished *IntValue  `xml:"yearpublished"`
	ProductCode   *ValueAttr `xml:"productcode"`
	// Width, Length, Depth and Height are in inches and Weight in pounds; 0 means not recorded
	Width  *FloatValue `xml:"width"`
	Length *FloatValue `xml:"length"`
	Depth  *FloatValue `xml:"depth"`
	Height *FloatValue `xml:"height"`
	Weight *FloatValue `xml:"weight"`
	// Format, PageCount, ISBN10, ISBN13 and Description are sent for RPG items
	Format      *ValueAttr `xml:"format"`
	PageCount   *IntValue  `xml:"pagecount"`
	ISBN10      *ValueAttr `xml:"isbn10"`
	ISBN13      *ValueAttr `xml:"isbn13"`
	Description string     `xml:"description"`
}

// Versions lists a thing's published versions
type Versions struct {
	XMLName xml.Name   `xml:"versions"`
	Items   []*Version `xml:"item"`
}

// Video is one video linked to a thing, returned with ThingsOptions.Videos
type Video struct {
	XMLName  xml.Name `xml:"video"`
	ID       int64    `xml:"id,attr"`
	Title    string   `xml:"title,attr"`
	Category string   `xml:"category,attr"`
	Language string   `xml:"language,attr"`
	// Link is the video's url, e.g. on YouTube
	Link     string    `xml:"link,attr"`
	Username string    `xml:"username,attr"`
	UserID   int64     `xml:"userid,attr"`
	PostDate time.Time `xml:"postdate,attr"`
}

// Videos lists a thing's videos
type Videos struct {
	XMLName xml.Name `xml:"videos"`
	Total   int64    `xml:"total,attr"`
	Videos  []*Video `xml:"video"`
}

// Comment is one user comment or rating on a thing
type Comment struct {
	XMLName  xml.Name `xml:"comment"`
	Username string   `xml:"username,attr"`
	Rating   string   `xml:"rating,attr"`
	Value    string   `xml:"value,attr"`
}

// Comments is one page of a thing's comments, returned with ThingsOptions.Comments or RatingComments
type Comments struct {
	XMLName    xml.Name   `xml:"comments"`
	Page       int        `xml:"page,attr"`
	TotalItems int64      `xml:"totalitems,attr"`
	Comments   []*Comment `xml:"comment"`
}

// Price is a GeekMarket asking price
type Price struct {
	XMLName  xml.Name `xml:"price"`
	Currency string   `xml:"currency,attr"`
	// Value is the amount as BGG sends it, e.g. "90.00"
	Value string `xml:"value,attr"`
}

// ListingLink links to a GeekMarket listing
type ListingLink struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
	Title   string   `xml:"title,attr"`
}

// MarketplaceListing is one GeekMarket listing for a thing, returned with ThingsOptions.MarketPlace
type MarketplaceListing struct {
	XMLName xml.Name `xml:"listing"`
	// ListDate is as BGG sends it, e.g. "Thu, 19 Jan 2006 22:08:15 +0000"
	ListDate  *ValueAttr   `xml:"listdate"`
	Price     *Price       `xml:"price"`
	Condition *ValueAttr   `xml:"condition"`
	Notes     *ValueAttr   `xml:"notes"`
	Link      *ListingLink `xml:"link"`
}

// MarketplaceListings lists a thing's GeekMarket listings
type MarketplaceListings struct {
	XMLName  xml.Name              `xml:"marketplacelistings"`
	Listings []*MarketplaceListing `xml:"listing"`
}

// Thing is a board game, expansion, accessory, RPG item or issue, or video game.
// Elements BGG does not send for the thing's type are nil or empty: player counts, play times and polls
// are sent for board games, ReleaseDate for video games, SeriesCode for RPG items, and DatePublished and
// IssueIndex for RPG issues.
type Thing struct {
	XMLName     xml.Name `xml:"item"`
	Type        string   `xml:"type,attr"`
	ID          int64    `xml:"id,attr"`
	Thumbnail   string   `xml:"thumbnail"`
	Image       string   `xml:"image"`
	Names       []*Name  `xml:"name"`
	Description string   `xml:"description"`

	YearPublished *IntValue    `xml:"yearpublished"`
	MinPlayers    *IntValue    `xml:"minplayers"`
	MaxPlayers    *IntValue    `xml:"maxplayers"`
	Polls         []*Poll      `xml:"poll"`
	PollSummary   *PollSummary `xml:"poll-summary"`
	// PlayingTime, MinPlayTime and MaxPlayTime are in minutes
	PlayingTime *IntValue `xml:"playingtime"`
	MinPlayTime *IntValue `xml:"minplaytime"`
	MaxPlayTime *IntValue `xml:"maxplaytime"`
	MinAge      *IntValue `xml:"minage"`

	// ReleaseDate is as BGG sends it and may have zero month or day, e.g. "2019-00-00"
	ReleaseDate   *ValueAttr `xml:"releasedate"`
	SeriesCode    *ValueAttr `xml:"seriescode"`
	DatePublished *ValueAttr `xml:"datepublished"`
	IssueIndex    *IntValue  `xml:"issueindex"`

	// Links include the thing's categories, mechanics, families, designers, artists, publishers and expansions
	Links []*Link `xml:"link"`

	Statistics          *Statistics          `xml:"statistics"`
	Versions            *Versions            `xml:"versions"`
	Videos              *Videos              `xml:"videos"`
	Comments            *Comments            `xml:"comments"`
	MarketplaceListings *MarketplaceListings `xml:"marketplacelistings"`
}

// ThingsResponse is the response from GetThings
type ThingsResponse struct {
	XMLName xml.Name `xml:"items"`
	Items   []*Thing `xml:"item"`
}
