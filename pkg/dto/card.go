package dto

type Cards struct {
	Cards []Card `json:"card"`
}

type Card struct {
	ID                    string   `json:"id,omitempty" bson:"id,omitempty"`
	Name                  string   `json:"name,omitempty" bson:"name,omitempty"`
	NationalPokedexNumber int      `json:"nationalPokedexNumber,omitempty" bson:"nationalPokedexNumber,omitempty"`
	ImageURL              string   `json:"imageUrl,omitempty" bson:"imageUrl,omitempty"`
	ImageURLHiRes         string   `json:"imageUrlHiRes,omitempty" bson:"imageUrlHiRes,omitempty"`
	Types                 []string `json:"types,omitempty" bson:"types,omitempty"`
	Supertype             string   `json:"supertype,omitempty" bson:"supertype,omitempty"`
	Subtype               string   `json:"subtype,omitempty"`
	Ability               struct {
		Name string `json:"name,omitempty"`
		Text string `json:"text,omitempty"`
		Type string `json:"type,omitempty"`
	} `json:"ability,omitempty"`
	AncientTrait struct {
		Name string `json:"name,omitempty"`
		Text string `json:"text,omitempty"`
	} `json:"ancientTrait,omitempty"`
	Hp                   interface{} `json:"hp,omitempty"`
	RetreatCost          []string    `json:"retreatCost,omitempty"`
	ConvertedRetreatCost int         `json:"convertedRetreatCost,omitempty"`
	Number               string      `json:"number,omitempty"`
	Artist               string      `json:"artist,omitempty"`
	Rarity               string      `json:"rarity,omitempty"`
	Series               string      `json:"series,omitempty"`
	Set                  string      `json:"set,omitempty"`
	SetCode              string      `json:"setCode,omitempty"`
	Text                 []string    `json:"text,omitempty"`
	Attacks              []struct {
		Cost                []string    `json:"cost,omitempty"`
		Name                string      `json:"name,omitempty"`
		Text                string      `json:"text,omitempty"`
		Damage              interface{} `json:"damage,omitempty"`
		ConvertedEnergyCost int         `json:"convertedEnergyCost,omitempty"`
	} `json:"attacks,omitempty"`
	Weaknesses []struct {
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"weaknesses,omitempty"`
	EvolvesFrom string `json:"evolvesFrom,omitempty"`
}
