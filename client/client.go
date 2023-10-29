// client/client.go
package client

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
)

// RestClient struct to manage REST client state
type RestClient struct {
	SessionID string
}

// Login performs a login operation and stores the sessionID.
func (c *RestClient) Login(username, password string) error {
	// Initialize a CookieJar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("Error initializing cookie jar: %v", err)
	}

	// Create an HTTP client with the cookie jar
	httpClient := &http.Client{
		Jar: jar,
	}

	// Prepare the multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.WriteField("user", username)
	writer.WriteField("password", password)
	writer.Close()

	// Create a new POST request
	req, err := http.NewRequest("POST", "https://rest.tastenext.de/public/login/process", &b)
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	// Set Content-Type for the request
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error performing request: %v", err)
	}
	defer resp.Body.Close()

	// Check for sessionid cookie
	for _, cookie := range jar.Cookies(req.URL) {
		fmt.Println(cookie.Name)
		if cookie.Name == "JSESSIONID" {
			c.SessionID = cookie.Value
			fmt.Println(cookie.Value)
			break
		}
	}
	return nil
}
