// client/client.go
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// RestClient struct to manage REST client state
type RestClient struct {
	Client    *http.Client
	SessionID string
	UserId    int
	CookieJar *cookiejar.Jar
}

func NewClient(config Config) (*RestClient, error) {
	c := &RestClient{}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing cookie jar: %v", err)
	}
	c.CookieJar = jar
	c.Client = &http.Client{
		Jar: jar,
	}
	cookie := &http.Cookie{
		Name:  "JSESSIONID",
		Value: config.SessionCredentials.SessionID,
		Path:  "/",
	}
	cookieUrl, err := url.Parse("https://rest.tastenext.de")
	c.CookieJar.SetCookies(cookieUrl, []*http.Cookie{cookie})

	c.SessionID = config.SessionCredentials.SessionID
	// Check if the old cookie works
	CurrentUserResponse, err := c.GetCurrentUser()
	// If not, login again
	if err != nil {
		fmt.Println("update session token")
		err := c.Login(config)
		if err != nil {
			return nil, fmt.Errorf("failed to log in")
		}
		// And check if it works now
		CurrentUserResponse, err = c.GetCurrentUser()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token")
		}
	}
	c.UserId = CurrentUserResponse.User.ID

	return c, nil
}

func (c *RestClient) SendRequest(method, urlStr string, body io.Reader, result interface{}) error {
	if c.Client == nil {
		return fmt.Errorf("client not initialized. Please login first")
	}

	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("server returns status code %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(respBody, result)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return nil
}

// Login performs a login operation and stores the sessionID.
func (c *RestClient) Login(config Config) error {
	// Prepare the multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.WriteField("username", config.LoginCredentials.User)
	writer.WriteField("password", config.LoginCredentials.Password)
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
	for _, cookie := range c.CookieJar.Cookies(req.URL) {
		if cookie.Name == "JSESSIONID" {
			c.SessionID = cookie.Value
			foundCookie = true
			break
		}
	}
	if !foundCookie {
		return fmt.Errorf("error getting session cookie")
	}

	return nil
}

// GetUser performs a GET request to get the current user
func (c *RestClient) GetCurrentUser() (CurrentUserResponse, error) {
	var currentUserResp CurrentUserResponse

	err := c.SendRequest("GET", "https://rest.tastenext.de/backend/user/current-user", nil, &currentUserResp)
	if err != nil {
		return CurrentUserResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	return currentUserResp, nil
}

func (c *RestClient) GetUser() (UserResponse, error) {
	var userResp UserResponse

	urlWithUserID := fmt.Sprintf("https://rest.tastenext.de/backend/user/%d", c.UserId)
	err := c.SendRequest("GET", urlWithUserID, nil, &userResp)
	if err != nil {
		log.Fatal("Error sending request")
	}
	return userResp, nil
}

func (c *RestClient) GetMenue() ([]MenuResponse, error) {
	var menuResponses []MenuResponse

	customer := 44897 // TODO: get this from the user object
	nextWeeks := getNextFourWeeks()

	for _, week := range nextWeeks {
		menuUrl := fmt.Sprintf(
			"https://rest.tastenext.de/frontend/menu/get-personal-menu-week/calendar-week/%d/year/%d/customer/%d/menu-block/14",
			week.CalendarWeek,
			week.Year,
			customer,
		)

		var menuResp MenuResponse
		err := c.SendRequest("GET", menuUrl, nil, &menuResp)
		if err != nil {
			log.Fatal("Error getting menues")
		}
		menuResponses = append(menuResponses, menuResp)
	}

	return menuResponses, nil
}
