package dto

import "fmt"

// Sets is a struct that contains an array of Set structs
type Sets struct {
	Sets []Set `json:"sets"`
}

// Set is a struct that holds all values regarding a specific Card set
type Set struct {
	Code          string `json:"code" bson:"code"`
	PtcgoCode     string `json:"ptcgoCode" bson:"ptcgoCode"`
	Name          string `json:"name" bson:"name"`
	Series        string `json:"series" bson:"series"`
	TotalCards    int    `json:"totalCards" bson:"totalCards"`
	StandardLegal bool   `json:"standardLegal" bson:"standardLegal"`
	ExpandedLegal bool   `json:"expandedLegal" bson:"expandedLegal"`
	ReleaseDate   string `json:"releaseDate" bson:"releaseDate"`
	SymbolURL     string `json:"symbolUrl" bson:"symbolUrl"`
	LogoURL       string `json:"logoUrl" bson:"logoUrl"`
	UpdatedAt     string `json:"updatedAt" bson:"updatedAt"`
}

// PrintSetNames will print the set names present in a Sets struct
func (sets *Sets) PrintSetNames() {

	fmt.Println("All Set Names: ")
	fmt.Println("------------------------")

	for _, s := range sets.Sets {
		fmt.Println(s.Name)
	}

	fmt.Println("------------------------")
}

// PrintStandardSets will print the set names present in a Sets struct whose StandardLegal field is "true"
func (sets *Sets) PrintStandardSets() {

	fmt.Println("Standard Sets: ")
	fmt.Println("------------------------")

	for _, s := range sets.Sets {

		if s.StandardLegal == true {
			fmt.Println(s.Name)
		}

	}

	fmt.Println("------------------------")
}
