package client

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
		Customer  struct {
			ID             int `json:"id"`
			AccountBalance struct {
				Amount int `json:"amount"`
			} `json:"accountBalance"`
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

type MenuResponse struct {
	Status      string      `json:"status"`
	Message     string      `json:"message"`
	Title       string      `json:"title"`
	Transmitted interface{} `json:"transmitted"`
	Bookings    []struct {
		ID                 int         `json:"id"`
		BookingPrice       int         `json:"bookingPrice"`
		BookingTime        int64       `json:"bookingTime"`
		PickupTime         interface{} `json:"pickupTime"`
		MenuBlockLineEntry struct {
			ID            int   `json:"id"`
			EmissionDate  int64 `json:"emissionDate"`
			MenuBlockLine struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Color       string `json:"color"`
				ColorBooked string `json:"colorBooked"`
			} `json:"menuBlockLine"`
			Dish struct {
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
			NumberOfBookings int `json:"numberOfBookings"`
		} `json:"menuBlockLineEntry"`
	} `json:"bookings"`
	Week                 int `json:"week"`
	Year                 int `json:"year"`
	MenuBlockWeekWrapper struct {
		MenuBlockWeek struct {
			MenuBlock struct {
				ID                              int    `json:"id"`
				Name                            string `json:"name"`
				BookingUntilTime                string `json:"bookingUntilTime"`
				CancellationUntilTime           string `json:"cancellationUntilTime"`
				BookingUntilXDaysInAdvance      int    `json:"bookingUntilXDaysInAdvance"`
				CancellationUntilXDaysInAdvance int    `json:"cancellationUntilXDaysInAdvance"`
			} `json:"menuBlock"`
			CalendarWeek             int `json:"calendarWeek"`
			Year                     int `json:"year"`
			NextCalendarWeek         int `json:"nextCalendarWeek"`
			NextCalendarWeekYear     int `json:"nextCalendarWeekYear"`
			PreviousCalendarWeek     int `json:"previousCalendarWeek"`
			PreviousCalendarWeekYear int `json:"previousCalendarWeekYear"`
			MenuBlockLineWeeks       []struct {
				MenuBlockLine struct {
					ID          int    `json:"id"`
					Name        string `json:"name"`
					Color       string `json:"color"`
					ColorBooked string `json:"colorBooked"`
				} `json:"menuBlockLine"`
				CalendarWeek int `json:"calendarWeek"`
				Year         int `json:"year"`
				Entries      []struct {
					ID            int   `json:"id"`
					EmissionDate  int64 `json:"emissionDate"`
					MenuBlockLine struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						Color       string `json:"color"`
						ColorBooked string `json:"colorBooked"`
					} `json:"menuBlockLine"`
					Dish struct {
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
					NumberOfBookings int `json:"numberOfBookings"`
				} `json:"entries"`
			} `json:"menuBlockLineWeeks"`
		} `json:"menuBlockWeek"`
		WeekdayHeadings []string `json:"weekdayHeadings"`
	} `json:"menuBlockWeekWrapper"`
	CurrentAccountBalance string      `json:"currentAccountBalance"`
	CustomerName          string      `json:"customerName"`
	CustomerID            int         `json:"customerId"`
	BpcJSON               interface{} `json:"bpcJson"`
}
