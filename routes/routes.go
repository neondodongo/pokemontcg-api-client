package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pokemontcg-api-client/dto"
)

// Immutable (unchangeable) variables
const (
	apiURI = "https://api.pokemontcg.io/v1" // URI to access the Pokemon TCG API
)

// GetAllSets calls the Pokemon TCG API to retrieve all existing sets
func GetAllSets() {

	// Send a Get request to the Pokemon TCG API to retrieve all existing set data
	// http.Get() returns a response pointer (specific memory address) and an error
	resp, err := http.Get(apiURI + "/sets")

	// If an error was returned from the Get request, print an error statement and leave the function at this point
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("Could not retrieve sets from API. Response from server: %d\n", resp.StatusCode)
		return
	}

	fmt.Printf("API call was a success! Response from server: %d\n", resp.StatusCode)

	// Defer will run after the rest of the function has finished
	// This specific defer closes the response body reader
	defer resp.Body.Close()

	// ioutil.ReadAll() will read the entire response Body from a successful Get request and returns a byte array and an error
	body, err := ioutil.ReadAll(resp.Body)

	// If an error was returned from the attempted reading of the response body received from the Pokemon TCG API, print an error statement and leave the function.
	if err != nil {
		fmt.Println("Error reading response body from Pokemon TCG API", err)
		return
	}

	// Declare a variable to hold the incoming Set data
	// See set.go in the dto package for specifics on the fields
	// This struct must comply with the payload received from the API
	var sets dto.Sets

	// Use the json module to unmarshal the response body to our Sets struct
	// json.Unmarshal returns an error if the Struct and the response body JSON do not match or the JSON is malformed
	// If an error is returned, print an error statement and exit the function
	fmt.Println("Unmarshal JSON to Sets struct...")
	if err := json.Unmarshal(body, &sets); err != nil {
		fmt.Println("Failed to unmarshal JSON to Sets struct. JSON may be malformed or is not compatible with the struct")
		return
	}

	// From this point, all of the set data should have been acquired and we can do whatever we want with it
	// This includes sending it elsewhere (to another web app), saving it to a database (SQL or MongoDB), or even creating files with the data
	// The Set data can also be used to get all Card data for each Set (using the SetCode field in another GET request to a different endpoint URI)

	fmt.Printf("Number of sets: %d\n", len(sets.Sets))
	fmt.Println("------------------------")

	//sets.PrintSetNames()
	sets.PrintStandardSets()
}
