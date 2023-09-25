package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getCurrentUser(cookieString string) (map[string]interface{}, error) {
	// Initialize HTTP client
	client := &http.Client{}

	// Create new HTTP request
	req, err := http.NewRequest("GET", "https://rest.tastenext.de/backend/user/current-user", nil)
	if err != nil {
		return nil, err
	}

	// Set cookie string in the request header
	req.Header.Set("Cookie", cookieString)

	// Execute the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch current user: %s", resp.Status)
	}

	// Read and parse the JSON response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user map[string]interface{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func main() {
	cookieString := "JSESSIONID=753551B645FB71F16A31ACAF7BFD310C"

	currentUser, err := getCurrentUser(cookieString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Successfully fetched current user: %+v\n", currentUser)
}
