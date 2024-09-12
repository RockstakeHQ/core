package types

type Goals struct {
	Home int `bson:"home" json:"home"`
	Away int `bson:"away" json:"away"`
}

type Score struct {
	Halftime  Halftime `bson:"halftime" json:"halftime"`
	Fulltime  Halftime `bson:"fulltime" json:"fulltime"`
	Extratime Halftime `bson:"extratime" json:"extratime"`
	Penalty   Halftime `bson:"penalty" json:"penalty"`
}

type Halftime struct {
	Home int `bson:"home" json:"home"`
	Away int `bson:"away" json:"away"`
}
