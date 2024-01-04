package client

import "time"

// Client types
type Dish struct {
	ID          int    `json:"ID"`
	Description string `json:"Description"`
	Name        string `json:"Name"`
}

type UpcomingDish struct {
	Dummy          bool
	Date           time.Time
	OrderId        int
	Dish           Dish
	Orders         int // We get the total order for each dish from the API ;-)
	PersonalOrders int // How often we ordered this dish in the past
	Booked         bool
}

type UpcomingDishMap map[string][]UpcomingDish

func (m UpcomingDishMap) Merge(other UpcomingDishMap) {
	for key, value := range other {
		m[key] = append(m[key], value...)
	}
}

type Bookings struct {
	ID                 int `json:"id"`
	BookingPrice       int `json:"bookingPrice"`
	MenuBlockLineEntry struct {
		ID   int  `json:"id"`
		Dish Dish `json:"dish"`
	} `json:"menuBlockLineEntry"`
}

type MenuDates struct {
	Year         int
	CalendarWeek int
}

// Response types
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
			Bookings []Bookings `json:"bookings"`
		} `json:"customer"`
	} `json:"user"`
}

type MenuResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Title    string `json:"title"`
	Bookings []struct {
		ID                 int   `json:"id"`
		BookingPrice       int   `json:"bookingPrice"`
		BookingTime        int64 `json:"bookingTime"`
		PickupTime         any   `json:"pickupTime"`
		MenuBlockLineEntry struct {
			ID               int  `json:"id"`
			Dish             Dish `json:"dish"`
			NumberOfBookings int  `json:"numberOfBookings"`
		} `json:"menuBlockLineEntry"`
	} `json:"bookings"`
	Week                 int `json:"week"`
	Year                 int `json:"year"`
	MenuBlockWeekWrapper struct {
		MenuBlockWeek struct {
			MenuBlock struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"menuBlock"`
			CalendarWeek       int `json:"calendarWeek"`
			Year               int `json:"year"`
			MenuBlockLineWeeks []struct {
				CalendarWeek int `json:"calendarWeek"`
				Year         int `json:"year"`
				Entries      []struct {
					EmissionDate     interface{} `json:"emissionDate"`
					ID               int         `json:"id"`
					Dish             Dish        `json:"dish"`
					NumberOfBookings int         `json:"numberOfBookings"`
				} `json:"entries"`
			} `json:"menuBlockLineWeeks"`
		} `json:"menuBlockWeek"`
	} `json:"menuBlockWeekWrapper"`
	CurrentAccountBalance string `json:"currentAccountBalance"`
	CustomerName          string `json:"customerName"`
	CustomerID            int    `json:"customerId"`
	BpcJSON               any    `json:"bpcJson"`
}
