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

type RestClient struct {
	Client     *http.Client
	Config     Config
	SessionID  string
	UserId     int
	CustomerId int
	CookieJar  *cookiejar.Jar
	Bookings   []Bookings
}

func NewClient(Config Config) (*RestClient, error) {
	c := &RestClient{}

	err := ValidateConfig(Config)
	if err != nil {
		return &RestClient{}, fmt.Errorf("validating config: %w", err)
	}
	c.Config = Config

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("initializing cookie jar: %v", err)
	}
	c.CookieJar = jar
	c.Client = &http.Client{
		Jar: jar,
	}
	cookie := &http.Cookie{
		Name:  "JSESSIONID",
		Value: c.Config.SessionCredentials.SessionID,
		Path:  "/",
	}
	cookieUrl, err := url.Parse("https://rest.tastenext.de")
	if err != nil {
		return &RestClient{}, fmt.Errorf("perse cookie url: %w", err)
	}
	c.CookieJar.SetCookies(cookieUrl, []*http.Cookie{cookie})

	// Check if the old SessionId works
	c.SessionID = c.Config.SessionCredentials.SessionID
	currentUserResponse, err := c.getCurrentUser()
	// If not, login again and get a new one
	if err != nil {
		fmt.Println("update session token")
		err := c.login()
		if err != nil {
			return nil, fmt.Errorf("failed to log in")
		}
		// Does it work now?
		currentUserResponse, err = c.getCurrentUser()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token")
		}
	}

	c.UserId = currentUserResponse.User.ID

	userResponse, err := c.GetUser()
	if err != nil {
		return &RestClient{}, fmt.Errorf("unable to load user")
	}
	c.CustomerId = userResponse.User.Customer.ID
	c.Bookings = userResponse.User.Customer.Bookings

	return c, nil
}

// Private
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

func (c *RestClient) login() error {
	// Prepare the multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.WriteField("username", c.Config.LoginCredentials.User)
	writer.WriteField("password", c.Config.LoginCredentials.Password)
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

	// Get sessionId from cookie
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

func (c *RestClient) getCurrentUser() (CurrentUserResponse, error) {
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

// We can only get one whole week from the API
func (c *RestClient) GetMenuWeek(Year int, Week int) (UpcomingDishMap, error) {
	var retVal = UpcomingDishMap{}

	customer := c.CustomerId

	menuUrl := fmt.Sprintf(
		"https://rest.tastenext.de/frontend/menu/get-personal-menu-week/calendar-week/%d/year/%d/customer/%d/menu-block/14",
		Week,
		Year,
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

			// Flag already booked dishes
			isBooked := false
			for _, booking := range menuResp.Bookings {
				if booking.MenuBlockLineEntry.ID == dish.ID {
					isBooked = true
				}
			}

			// Append upcoming dishes
			personalOrderCount, _ := c.GetOrderCount(dish.Dish.ID)
			upcomingDish := UpcomingDish{
				OrderId:        dish.ID,
				Dish:           dish.Dish,
				Orders:         dish.NumberOfBookings,
				PersonalOrders: personalOrderCount,
				Date:           edate,
				Dummy:          isDummy,
				Booked:         isBooked,
			}
			dateKey := edate.Format("06-01-02")
			retVal[dateKey] = append(retVal[dateKey], upcomingDish)
		}
	}
	return retVal, nil
}

// Get menu for the next n calender weeks
func (c *RestClient) GetMenuWeeks(weeks int) (UpcomingDishMap, error) {
	var retVal = UpcomingDishMap{}

	nextWeeks := GetNextCalenderWeeks(weeks)

	for _, week := range nextWeeks {
		menuWeek, err := c.GetMenuWeek(week.Year, week.CalendarWeek)
		if err != nil {
			fmt.Errorf("Error getting weeks")
		}
		retVal.Merge(menuWeek)
	}

	return retVal, nil
}

// Get Menu for one Day
func (c *RestClient) GetMenuDay(Day time.Time) (UpcomingDishMap, error) {
	var retVal = UpcomingDishMap{}

	menuWeek, err := c.GetMenuWeek(Day.ISOWeek())
	if err != nil {
		return retVal, fmt.Errorf("error: %w", err)
	}

	dateKey := Day.Format("06-01-02")
	retVal[dateKey] = menuWeek[dateKey]
	if len(retVal) == 0 {
		return retVal, fmt.Errorf("no dishes found for this day")
	}
	return retVal, nil
}

// Order or cancel a dish
func (c *RestClient) OrderDish(DishOrderId int, CancelOrder bool) error {
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

	// Toggle booking
	bookingUrl := fmt.Sprintf(
		"https://rest.tastenext.de/frontend/menu/order/menu-block-line-entry/%d/customer/%d",
		DishOrderId,
		c.CustomerId)

	var menuResp MenuResponse
	err = c.sendRequest("GET", bookingUrl, nil, &menuResp)
	if err != nil {
		return errors.New("failed sending order request")
	}

	switch menuResp.Message {
	case "app.messages.changed-booking-status.too-late":
		return fmt.Errorf("to late to place order")

	case "app.messages.changed-booking-status.insufficient-money":
		return fmt.Errorf("not enough account balance to place order")

	case "Ok":
		return nil

	default:
		return fmt.Errorf("failed to place / remove order: %v", menuResp.Message)
	}
}

// How often a dish was ordered in the past
func (c RestClient) GetOrderCount(DishId int) (count int, dish Dish) {
	count = 0
	dish = Dish{}

	for _, booking := range c.Bookings {
		if DishId == booking.MenuBlockLineEntry.Dish.ID {
			count++
			dish = booking.MenuBlockLineEntry.Dish
		}
	}

	return count, dish
}
