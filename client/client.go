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

// CurrentUserResponse struct for unmarshaling the user ID
type CurrentUserResponse struct {
	User struct {
		ID int `json:"id"`
	} `json:"user"`
}

type UserResponse struct {
	User struct {
		ID              int    `json:"id"`
		Email           string `json:"email"`
		FirstName       string `json:"firstName"`
		InitialPassword string `json:"initialPassword"`
		LastName        string `json:"lastName"`
		OldID           int    `json:"oldId"`
		Phone           string `json:"phone"`
		RfidKey         string `json:"rfidKey"`
		Username        string `json:"username"`
		Locked          bool   `json:"locked"`
		Active          bool   `json:"active"`
		CatererAdmin    any    `json:"catererAdmin"`
		Customer        struct {
			ID             int `json:"id"`
			AccountBalance struct {
				Amount int `json:"amount"`
			} `json:"accountBalance"`
			Client struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				NameShort string `json:"nameShort"`
				Caterer   struct {
					ID                    int    `json:"id"`
					Name                  string `json:"name"`
					Recipient             string `json:"recipient"`
					Iban                  string `json:"iban"`
					Bic                   any    `json:"bic"`
					TopupReasonForPayment string `json:"topupReasonForPayment"`
				} `json:"caterer"`
			} `json:"client"`
			DefaultPriceLevel struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"defaultPriceLevel"`
			Payments []struct {
				ID     int `json:"id"`
				Amount struct {
					Amount int `json:"amount"`
				} `json:"amount"`
				PaymentMethod     string `json:"paymentMethod"`
				Currency          any    `json:"currency"`
				Iban              any    `json:"iban"`
				Bic               any    `json:"bic"`
				Payer             any    `json:"payer"`
				ReasonForPayment  any    `json:"reasonForPayment"`
				BankTransactionID any    `json:"bankTransactionId"`
				TimeOfBooking     int64  `json:"timeOfBooking"`
				BookkeepingDate   any    `json:"bookkeepingDate"`
				ValueDate         any    `json:"valueDate"`
			} `json:"payments"`
			Bookings []struct {
				ID                 int   `json:"id"`
				BookingPrice       int   `json:"bookingPrice"`
				BookingTime        int64 `json:"bookingTime"`
				PickupTime         any   `json:"pickupTime"`
				MenuBlockLineEntry struct {
					ID           int   `json:"id"`
					EmissionDate int64 `json:"emissionDate"`
					Dish         struct {
						ID          int    `json:"id"`
						Description string `json:"description"`
						Name        string `json:"name"`
						Additives   []struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Identifier string `json:"identifier"`
							Type       string `json:"type"`
						} `json:"additives"`
					} `json:"dish"`
				} `json:"menuBlockLineEntry"`
			} `json:"bookings"`
			Attributes struct {
			} `json:"attributes"`
			CreatedByUsername      string `json:"createdByUsername"`
			LastModifiedByUsername string `json:"lastModifiedByUsername"`
		} `json:"customer"`
		SystemAdmin    any `json:"systemAdmin"`
		ClassLevel     any `json:"classLevel"`
		ClassQualifier any `json:"classQualifier"`
	} `json:"user"`
	AvailablePriceLevels []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"availablePriceLevels"`
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

	CurrentUserResponse, err := c.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("error getting current user")
	}

	c.UserId = CurrentUserResponse.User.ID
	return nil
}

// GetUser performs a GET request to get the current user
func (c *RestClient) GetCurrentUser() (CurrentUserResponse, error) {
	var currentUserResp CurrentUserResponse

	if c.Client == nil {
		return currentUserResp, fmt.Errorf("client not initialized. Please login first")
	}

	req, err := http.NewRequest("GET", "https://rest.tastenext.de/backend/user/current-user", nil)
	if err != nil {
		return currentUserResp, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return currentUserResp, fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return currentUserResp, fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(body, &currentUserResp)
	if err != nil {
		return currentUserResp, fmt.Errorf("error unmarshaling JSON: %v", err)
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

	// Pretty-print the JSON
	prettyJSON, err := json.MarshalIndent(userResp, "", "  ")
	if err != nil {
		return userResp, fmt.Errorf("error marshaling JSON: %v", err)
	}

	fmt.Println(string(prettyJSON))

	return userResp, nil
}
