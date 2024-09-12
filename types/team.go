package types

type Team struct {
	ID     int    `bson:"_id" json:"id"`
	Name   string `bson:"name" json:"name"`
	Logo   string `bson:"logo" json:"logo"`
	Winner bool   `bson:"winner" json:"winner"`
}

type Teams struct {
	Home Team `bson:"home" json:"home"`
	Away Team `bson:"away" json:"away"`
}
