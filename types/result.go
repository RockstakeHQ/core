package types

type Goals struct {
	Home int `bun:"home" json:"home"`
	Away int `bun:"away" json:"away"`
}

type Score struct {
	Halftime  Halftime `bun:"halftime" json:"halftime"`
	Fulltime  Halftime `bun:"fulltime" json:"fulltime"`
	Extratime Halftime `bun:"extratime" json:"extratime"`
	Penalty   Halftime `bun:"penalty" json:"penalty"`
}

type Halftime struct {
	Home int `bun:"home" json:"home"`
	Away int `bun:"away" json:"away"`
}
