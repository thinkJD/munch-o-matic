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

	if c.Client == nil {
		return CurrentUserResponse{}, fmt.Errorf("client not initialized. Please login first")
	}

	req, err := http.NewRequest("GET", "https://rest.tastenext.de/backend/user/current-user", nil)
	if err != nil {
		return CurrentUserResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return CurrentUserResponse{}, fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return CurrentUserResponse{}, fmt.Errorf("Server returned status code %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CurrentUserResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(body, &currentUserResp)
	if err != nil {
		return CurrentUserResponse{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return currentUserResp, nil
}

func (c *RestClient) GetUser() (UserResponse, error) {
	var userResp UserResponse

	if c.Client == nil {
		return userResp, fmt.Errorf("client not initialized. Please login first")
	}

	urlWithUserID := fmt.Sprintf("https://rest.tastenext.de/backend/user/%d", c.UserId)

	req, err := http.NewRequest("GET", urlWithUserID, nil)
	if err != nil {
		return UserResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return UserResponse{}, fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(body, &userResp)
	if err != nil {
		return UserResponse{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Pretty-print the JSON
	prettyJSON, err := json.MarshalIndent(userResp, "", "  ")
	if err != nil {
		return UserResponse{}, fmt.Errorf("error marshaling JSON: %v", err)
	}

	fmt.Println(string(prettyJSON))

	return userResp, nil
}

func (c *RestClient) GetMenue() (MenuResponse, error) {
	var menuResponse MenuResponse

	if c.Client == nil {
		return MenuResponse{}, fmt.Errorf("client not initialized. Please login first")
	}
	customer := 44897 // TODO: get this from the user object

	nextWeeks := getNextFourWeeks()

	menuUrl := fmt.Sprintf(
		"https://rest.tastenext.de/frontend/menu/get-personal-menu-week/calendar-week/%d/year/%d/customer/%d/menu-block/14",
		nextWeeks[0].CalendarWeek,
		nextWeeks[0].Year,
		customer,
	)

	fmt.Println(menuUrl)

	return menuResponse, nil
}
