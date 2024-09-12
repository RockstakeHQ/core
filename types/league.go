package types

type League struct {
	ID      int    `bson:"_id" json:"id"`
	Name    string `bson:"name" json:"name"`
	Country string `bson:"country" json:"country"`
	Logo    string `bson:"logo" json:"logo"`
	Flag    string `bson:"flag" json:"flag"`
}
