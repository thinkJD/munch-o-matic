package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Add more fields as needed
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

func (c *Client) Login(username, password string) error {
	// Create a request body with the provided username and password
	reqBody := map[string]string{
		"username": username,
		"password": password,
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Make a POST request to the login endpoint
	loginURL := fmt.Sprintf("%s/login", c.BaseURL)
	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle the response and return any errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status code %d", resp.StatusCode)
	}

	// Process the response body if needed
	// ...

	return nil
}

func (c *Client) GetUser(userID int) (*User, error) {
	// Make a GET request to the user endpoint with the provided userID
	userURL := fmt.Sprintf("%s/user/%d", c.BaseURL, userID)
	resp, err := http.Get(userURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle the response and return any errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user with status code %d", resp.StatusCode)
	}

	// Parse the response into a User object
	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
