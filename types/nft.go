package types

type NftNodeInfo struct {
	FixtureID string  `json:"fixture_id"`
	Market    string  `json:"market"`
	Selection string  `json:"selection"`
	Type      string  `json:"type"`
	Odd       float64 `json:"odd"`
	Stake     float64 `json:"stake"`
}

type NFTMetadata struct {
	Description string      `json:"description"`
	Attributes  []Attribute `json:"attributes"`
	Compiler    string      `json:"compiler"`
}

type Attribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}
