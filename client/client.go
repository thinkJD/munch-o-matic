// client/client.go
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
)

// RestClient struct to manage REST client state
type RestClient struct {
	Client    *http.Client
	SessionID string
	UserId    int
}

// UserResponse struct for unmarshaling the user ID
type UserResponse struct {
	User struct {
		ID int `json:"id"`
	} `json:"user"`
}

// Login performs a login operation and stores the sessionID.
func (c *RestClient) Login(username, password string) error {
	// Initialize a CookieJar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("error initializing cookie jar: %v", err)
	}

	// Create an HTTP client with the cookie jar
	c.Client = &http.Client{
		Jar: jar,
	}

	// Prepare the multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.WriteField("username", username)
	writer.WriteField("password", password)
	writer.WriteField("remember-me", "true")
	writer.Close()

	// Create a new POST request
	req, err := http.NewRequest("POST", "https://rest.tastenext.de/public/login/process", &b)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set Content-Type for the request
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	foundCookie := false
	for _, cookie := range jar.Cookies(req.URL) {
		if cookie.Name == "JSESSIONID" {
			c.SessionID = cookie.Value
			foundCookie = true
			break
		}
	}
	if !foundCookie {
		return fmt.Errorf("error getting session cookie")
	}

	UserResponse, err := c.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("error getting current user")
	}

	c.UserId = UserResponse.User.ID
	fmt.Println("%s", c.UserId)
	return nil
}

// GetUser performs a GET request to get the current user
func (c *RestClient) GetCurrentUser() (UserResponse, error) {
	var userResp UserResponse

	if c.Client == nil {
		return userResp, fmt.Errorf("client not initialized. Please login first")
	}

	req, err := http.NewRequest("GET", "https://rest.tastenext.de/backend/user/current-user", nil)
	if err != nil {
		return userResp, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return userResp, fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return userResp, fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(body, &userResp)
	if err != nil {
		return userResp, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return userResp, nil
}
