package client

// CurrentUserResponse struct for unmarshaling the user ID
type CurrentUserResponse struct {
	User struct {
		ID int `json:"id"`
	} `json:"user"`
}

type UserResponse struct {
	User struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		OldID     int    `json:"oldId"`
		Phone     string `json:"phone"`
		RfidKey   string `json:"rfidKey"`
		Username  string `json:"username"`
		Locked    bool   `json:"locked"`
		Active    bool   `json:"active"`
		Customer  struct {
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
		} `json:"customer"`
	} `json:"user"`
}
