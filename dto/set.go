package dto

import "fmt"

// Sets is a struct that contains an array of Set structs
type Sets struct {
	Sets []Set `json:"sets"`
}

// Set is a struct that holds all values regarding a specific Card set
type Set struct {
	Code          string `json:"code"`
	PtcgoCode     string `json:"ptcgoCode"`
	Name          string `json:"name"`
	Series        string `json:"series"`
	TotalCards    int    `json:"totalCards"`
	StandardLegal bool   `json:"standardLegal"`
	ExpandedLegal bool   `json:"expandedLegal"`
	ReleaseDate   string `json:"releaseDate"`
	SymbolURL     string `json:"symbolUrl"`
	LogoURL       string `json:"logoUrl"`
	UpdatedAt     string `json:"updatedAt"`
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
