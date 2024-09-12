package types

type BetType struct {
	ID   int    `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type BetValue struct {
	ID    int    `bson:"_id" json:"id"`
	Value string `bson:"value" json:"value"`
}
