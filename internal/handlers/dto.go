package handlers

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *UserResponse `json:"user"`
}

type UserResponse struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	PassportData string `json:"passportData,omitempty"`
	Role         string `json:"role"`
}

type UpdateUserRequest struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	PassportData string `json:"passportData"`
}

type RouteSearchResponse struct {
	RouteID       int64   `json:"routeId"`
	TrainNumber   string  `json:"trainNumber"`
	DepartureTime string  `json:"departureTime"`
	ArrivalTime   string  `json:"arrivalTime"`
	Price         float64 `json:"price"`
	AvailableSeats int    `json:"availableSeats"`
}

type CreateOrderRequest struct {
	RouteID     int64   `json:"routeId" binding:"required"`
	SeatID      int64   `json:"seatId" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	PassengerID *int64  `json:"passengerId"`
}

type OrderResponse struct {
	ID         int64          `json:"id"`
	UserID     int64          `json:"userId"`
	RouteID    *int64         `json:"routeId,omitempty"`
	RouteName  string         `json:"routeName,omitempty"`
	TrainNumber string        `json:"trainNumber,omitempty"`
	TrainType  string         `json:"trainType,omitempty"`
	DepartureCity string      `json:"departureCity,omitempty"`
	ArrivalCity string        `json:"arrivalCity,omitempty"`
	DepartureTime string      `json:"departureTime,omitempty"`
	ArrivalTime string        `json:"arrivalTime,omitempty"`
	CreatedAt  string         `json:"createdAt"`
	Status     string         `json:"status"`
	TotalAmount float64       `json:"totalAmount"`
	Tickets    []TicketResponse `json:"tickets"`
}

type TicketResponse struct {
	ID           int64   `json:"id"`
	TicketNumber string  `json:"ticketNumber"`
	SeatNumber   *int    `json:"seatNumber"`
	CarriageNumber *int   `json:"carriageNumber"`
	Price        float64 `json:"price"`
}

type PaymentRequest struct {
	CardNumber string `json:"cardNumber" binding:"required"`
	ExpiryDate string `json:"expiryDate" binding:"required"`
	CVV        string `json:"cvv" binding:"required"`
}

type PaymentResponse struct {
	Status        string `json:"status"`
	TransactionID string `json:"transactionId"`
}

type CreateRouteRequest struct {
	Name    string `json:"name" binding:"required"`
	TrainID int64  `json:"trainId" binding:"required"`
}

type UpdateRouteRequest struct {
	Name    string `json:"name"`
	TrainID int64  `json:"trainId"`
}

type CreateTrainRequest struct {
	Number string `json:"number" binding:"required"`
	Type   string `json:"type"`
}

type UpdateTrainRequest struct {
	Number string `json:"number"`
	Type   string `json:"type"`
}

