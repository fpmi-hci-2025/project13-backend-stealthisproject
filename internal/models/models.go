package models

import "time"

type User struct {
	ID           int64     `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type Passenger struct {
	ID           int64  `json:"id" db:"id"`
	UserID       int64  `json:"userId" db:"user_id"`
	FirstName    string `json:"firstName" db:"first_name"`
	LastName     string `json:"lastName" db:"last_name"`
	PassportData string `json:"passportData" db:"passport_data"`
}

type Train struct {
	ID     int64  `json:"id" db:"id"`
	Number string `json:"number" db:"number"`
	Type   string `json:"type" db:"type"`
}

type Carriage struct {
	ID      int64  `json:"id" db:"id"`
	TrainID int64  `json:"trainId" db:"train_id"`
	Number  int    `json:"number" db:"number"`
	Type    string `json:"type" db:"type"`
}

type Seat struct {
	ID         int64 `json:"id" db:"id"`
	CarriageID int64 `json:"carriageId" db:"carriage_id"`
	Number     int   `json:"number" db:"number"`
}

type Station struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	City string `json:"city" db:"city"`
}

type Route struct {
	ID     int64   `json:"id" db:"id"`
	Name   string  `json:"name" db:"name"`
	TrainID int64  `json:"trainId" db:"train_id"`
	Price  float64 `json:"price" db:"price"`
}

type RouteStation struct {
	RouteID      int64     `json:"routeId" db:"route_id"`
	StationID    int64     `json:"stationId" db:"station_id"`
	ArrivalTime  *time.Time `json:"arrivalTime" db:"arrival_time"`
	DepartureTime *time.Time `json:"departureTime" db:"departure_time"`
	StopOrder    int       `json:"stopOrder" db:"stop_order"`
}

type Order struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"userId" db:"user_id"`
	RouteID    *int64    `json:"routeId,omitempty" db:"route_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	Status     string    `json:"status" db:"status"`
	TotalAmount float64  `json:"totalAmount" db:"total_amount"`
	Tickets    []Ticket  `json:"tickets,omitempty"`
}

type Ticket struct {
	ID           int64     `json:"id" db:"id"`
	OrderID      int64     `json:"orderId" db:"order_id"`
	SeatID       *int64    `json:"seatId" db:"seat_id"`
	PassengerID  *int64    `json:"passengerId" db:"passenger_id"`
	DepartureDate time.Time `json:"departureDate" db:"departure_date"`
	Price        float64   `json:"price" db:"price"`
	TicketNumber string    `json:"ticketNumber" db:"ticket_number"`
	Status       string    `json:"status" db:"status"`
}

