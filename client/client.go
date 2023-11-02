// client/client.go
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// RestClient struct to manage REST client state
type RestClient struct {
	Client     *http.Client
	SessionID  string
	UserId     int
	CustomerId int
	CookieJar  *cookiejar.Jar
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

	// Check if the old SessionId works
	c.SessionID = config.SessionCredentials.SessionID
	currentUserResponse, err := c.GetCurrentUser()
	// If not, login again and get a new one
	if err != nil {
		fmt.Println("update session token")
		err := c.Login(config)
		if err != nil {
			return nil, fmt.Errorf("failed to log in")
		}
		// Does it work now?
		currentUserResponse, err = c.GetCurrentUser()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token")
		}
	}

	c.UserId = currentUserResponse.User.ID

	userResponse, err := c.GetUser()
	if err != nil {
		return &RestClient{}, fmt.Errorf("unable to load user", err)
	}
	c.CustomerId = userResponse.User.Customer.ID

	return c, nil
}

func (c *RestClient) sendRequest(method, urlStr string, body io.Reader, result interface{}) error {
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

	err := c.sendRequest("GET", "https://rest.tastenext.de/backend/user/current-user", nil, &currentUserResp)
	if err != nil {
		return CurrentUserResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	return currentUserResp, nil
}

func (c *RestClient) GetUser() (UserResponse, error) {
	var userResp UserResponse

	urlWithUserID := fmt.Sprintf("https://rest.tastenext.de/backend/user/%d", c.UserId)
	err := c.sendRequest("GET", urlWithUserID, nil, &userResp)
	if err != nil {
		log.Fatal("Error sending request")
	}
	return userResp, nil
}

type UpcomingDish struct {
	Dummy   bool
	Date    time.Time
	OrderId int
	Dish    Dish
	Orders  int // We get the total order for each dish from the API ;-)
	Booked  bool
}

func (c *RestClient) GetMenu() (map[string][]UpcomingDish, error) {
	var upcomingDishes = map[string][]UpcomingDish{}

	customer := c.CustomerId
	nextWeeks := getNextFourWeeks()

	for _, week := range nextWeeks {
		menuUrl := fmt.Sprintf(
			"https://rest.tastenext.de/frontend/menu/get-personal-menu-week/calendar-week/%d/year/%d/customer/%d/menu-block/14",
			week.CalendarWeek,
			week.Year,
			customer,
		)

		var menuResp MenuResponse
		err := c.sendRequest("GET", menuUrl, nil, &menuResp)
		if err != nil {
			log.Fatal("Error getting menus")
		}

		// extract fields
		for _, mblw := range menuResp.MenuBlockWeekWrapper.MenuBlockWeek.MenuBlockLineWeeks {
			for _, dish := range mblw.Entries {
				edate, err := GetEmissionDateAsTime(dish.EmissionDate)
				if err != nil {
					log.Fatal("Error getting emission date")
				}

				// Check for dummy values. They appear if there is no menu for that day.
				isDummy := dish.Dish.Name == "---"

				// Check bookings for this week
				isBooked := false
				for _, booking := range menuResp.Bookings {
					if booking.MenuBlockLineEntry.ID == dish.ID {
						isBooked = true
					}
				}

				upcomingDish := UpcomingDish{
					OrderId: dish.ID,
					Dish:    dish.Dish,
					Orders:  dish.NumberOfBookings,
					Date:    edate,
					Dummy:   isDummy,
					Booked:  isBooked,
				}
				dateKey := edate.Format("06-01-02")
				upcomingDishes[dateKey] = append(upcomingDishes[dateKey], upcomingDish)
			}
		}
	}

	return upcomingDishes, nil
}

func (c *RestClient) OrderMenu(DishOrderId int, CancelOrder bool) error {
	// Is the dish already ordered?
	userResp, err := c.GetUser()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	var alreadyOrdered = false
	for _, booking := range userResp.User.Customer.Bookings {
		if booking.MenuBlockLineEntry.ID == DishOrderId {
			alreadyOrdered = true
			break
		}
	}

	// Check if there is something to do, return if not
	if (alreadyOrdered && !CancelOrder) || (!alreadyOrdered && CancelOrder) {
		return nil
	}

	// toggle order
	bookingUrl := fmt.Sprintf(
		"https://rest.tastenext.de/frontend/menu/order/menu-block-line-entry/%d/customer/%d",
		DishOrderId,
		c.CustomerId)

	var menuResp MenuResponse
	err = c.sendRequest("GET", bookingUrl, nil, &menuResp)
	if err != nil {
		return errors.New("failed sending order request")
	}

	if menuResp.Status != "OK" {
		return fmt.Errorf("failed to place / remove order: %v", menuResp.Message)
	}

	return nil
}
