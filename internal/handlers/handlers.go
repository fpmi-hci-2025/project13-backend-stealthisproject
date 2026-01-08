package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/project13/backend-stealthisproject/internal/models"
	"github.com/project13/backend-stealthisproject/internal/repository"
	"github.com/project13/backend-stealthisproject/pkg/auth"
)

type Handlers struct {
	repos       *repository.Repositories
	authService *auth.AuthService
}

func NewHandlers(repos *repository.Repositories, authService *auth.AuthService) *Handlers {
	return &Handlers{
		repos:       repos,
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new passenger account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *Handlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	existingUser, _ := h.repos.User.GetByEmail(req.Email)
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	}

	// Hash password
	passwordHash, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         "PASSENGER",
	}
	if err := h.repos.User.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create passenger profile
	passenger := &models.Passenger{
		UserID:    user.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	if err := h.repos.Passenger.Create(passenger); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create passenger profile"})
		return
	}

	// Generate token
	token, err := h.authService.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User: &UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Role:      user.Role,
		},
	})
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repos.User.GetByEmail(req.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := h.authService.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetCurrentUser returns current user profile
// @Summary Get current user
// @Description Get authenticated user's profile
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 401 {object} map[string]string
// @Router /users/me [get]
func (h *Handlers) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int64)

	user, err := h.repos.User.GetByID(id)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	passenger, _ := h.repos.Passenger.GetByUserID(id)
	response := UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	if passenger != nil {
		response.FirstName = passenger.FirstName
		response.LastName = passenger.LastName
		response.PassportData = passenger.PassportData
	}

	c.JSON(http.StatusOK, response)
}

// UpdateCurrentUser updates current user profile
// @Summary Update current user
// @Description Update authenticated user's profile
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body UpdateUserRequest true "Update data"
// @Success 200 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Router /users/me [put]
func (h *Handlers) UpdateCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int64)

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passenger, err := h.repos.Passenger.GetByUserID(id)
	if err != nil || passenger == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Passenger profile not found"})
		return
	}

	if req.FirstName != "" {
		passenger.FirstName = req.FirstName
	}
	if req.LastName != "" {
		passenger.LastName = req.LastName
	}
	// Allow empty string to clear passport data
	passenger.PassportData = req.PassportData

	if err := h.repos.Passenger.Update(passenger); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	user, _ := h.repos.User.GetByID(id)
	response := UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    passenger.FirstName,
		LastName:     passenger.LastName,
		PassportData: passenger.PassportData,
		Role:         user.Role,
	}

	c.JSON(http.StatusOK, response)
}

// SearchRoutes searches for routes
// @Summary Search routes
// @Description Search for routes between cities
// @Tags Routes
// @Accept json
// @Produce json
// @Param from_city query string true "Departure city"
// @Param to_city query string true "Arrival city"
// @Param date query string true "Travel date"
// @Success 200 {array} RouteSearchResponse
// @Router /routes/search [get]
func (h *Handlers) SearchRoutes(c *gin.Context) {
	fromCity := c.Query("from_city")
	toCity := c.Query("to_city")
	date := c.Query("date")

	if fromCity == "" || toCity == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from_city, to_city, and date are required"})
		return
	}

	routes, err := h.repos.Route.Search(fromCity, toCity, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to search routes: %v", err)})
		return
	}

	var responses []RouteSearchResponse
	for _, route := range routes {
		train, _ := h.repos.Train.GetByID(route.TrainID)
		routeStations, _ := h.repos.Route.GetStations(route.ID)

		var departureTime, arrivalTime string
		if len(routeStations) > 0 {
			if routeStations[0].DepartureTime != nil {
				// Format time as ISO 8601 with date (use today's date + time from DB)
				depTime := *routeStations[0].DepartureTime
				// Combine today's date with the time from database
				today := time.Now()
				combinedTime := time.Date(today.Year(), today.Month(), today.Day(),
					depTime.Hour(), depTime.Minute(), depTime.Second(), depTime.Nanosecond(), time.UTC)
				departureTime = combinedTime.Format(time.RFC3339)
			}
			if len(routeStations) > 1 && routeStations[len(routeStations)-1].ArrivalTime != nil {
				arrTime := *routeStations[len(routeStations)-1].ArrivalTime
				today := time.Now()
				combinedTime := time.Date(today.Year(), today.Month(), today.Day(),
					arrTime.Hour(), arrTime.Minute(), arrTime.Second(), arrTime.Nanosecond(), time.UTC)
				arrivalTime = combinedTime.Format(time.RFC3339)
			}
		}

		// Calculate available seats (simplified)
		carriages, _ := h.repos.Carriage.GetByTrainID(route.TrainID)
		availableSeats := 0
		for _, carriage := range carriages {
			seats, _ := h.repos.Seat.GetByCarriageID(carriage.ID)
			availableSeats += len(seats)
		}

		trainNumber := ""
		if train != nil {
			trainNumber = train.Number
		}

		responses = append(responses, RouteSearchResponse{
			RouteID:       route.ID,
			TrainNumber:   trainNumber,
			DepartureTime: departureTime,
			ArrivalTime:   arrivalTime,
			Price:         route.Price, // Use actual route price from database
			AvailableSeats: availableSeats,
		})
	}

	// Always return an array, even if empty
	if responses == nil {
		responses = []RouteSearchResponse{}
	}
	c.JSON(http.StatusOK, responses)
}

// GetRoute gets route details
// @Summary Get route details
// @Description Get detailed information about a route
// @Tags Routes
// @Produce json
// @Param id path int true "Route ID"
// @Success 200 {object} map[string]interface{}
// @Router /routes/{id} [get]
func (h *Handlers) GetRoute(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID"})
		return
	}

	route, err := h.repos.Route.GetByID(id)
	if err != nil || route == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	train, _ := h.repos.Train.GetByID(route.TrainID)
	routeStations, _ := h.repos.Route.GetStations(route.ID)

	var stations []map[string]interface{}
	for _, rs := range routeStations {
		station, _ := h.repos.Station.GetByID(rs.StationID)
		stations = append(stations, map[string]interface{}{
			"station":       station,
			"arrivalTime":   rs.ArrivalTime,
			"departureTime": rs.DepartureTime,
			"stopOrder":     rs.StopOrder,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       route.ID,
		"name":     route.Name,
		"train":    train,
		"stations": stations,
	})
}

// CreateOrder creates a new order
// @Summary Create order
// @Description Create a new ticket order
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateOrderRequest true "Order data"
// @Success 201 {object} OrderResponse
// @Failure 400 {object} map[string]string
// @Router /orders [post]
func (h *Handlers) CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int64)

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check seat availability
	seat, err := h.repos.Seat.GetByID(req.SeatID)
	if err != nil || seat == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seat not found"})
		return
	}

	// Get route to find train
	route, err := h.repos.Route.GetByID(req.RouteID)
	if err != nil || route == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Route not found"})
		return
	}

	// Create order with price from request
	routeID := route.ID
	order := &models.Order{
		UserID:     id,
		RouteID:    &routeID,
		Status:     "PENDING",
		TotalAmount: req.Price,
	}
	if err := h.repos.Order.Create(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Determine passenger ID
	passengerID := req.PassengerID
	if passengerID == nil {
		passenger, _ := h.repos.Passenger.GetByUserID(id)
		if passenger != nil {
			pid := passenger.ID
			passengerID = &pid
		}
	}

	// Create ticket with price from request
	ticketNumber := fmt.Sprintf("TK-%d-%d", order.ID, time.Now().Unix())
	ticket := &models.Ticket{
		OrderID:      order.ID,
		SeatID:       &req.SeatID,
		PassengerID:  passengerID,
		DepartureDate: time.Now().AddDate(0, 0, 1), // Tomorrow
		Price:        req.Price,
		TicketNumber: ticketNumber,
		Status:       "ACTIVE",
	}
	if err := h.repos.Ticket.Create(ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	// Get seat and carriage info
	carriage, _ := h.repos.Carriage.GetByID(seat.CarriageID)

	// Get route information
	var routeName, trainNumber, trainType, departureCity, arrivalCity, departureTime, arrivalTime string
	if route != nil {
		routeName = route.Name
		train, _ := h.repos.Train.GetByID(route.TrainID)
		if train != nil {
			trainNumber = train.Number
			trainType = train.Type
		}
		routeStations, _ := h.repos.Route.GetStations(route.ID)
		if len(routeStations) > 0 {
			departureStation, _ := h.repos.Station.GetByID(routeStations[0].StationID)
			if departureStation != nil {
				departureCity = departureStation.City
			}
			if routeStations[0].DepartureTime != nil {
				departureTime = routeStations[0].DepartureTime.Format("15:04")
			}
			if len(routeStations) > 1 {
				arrivalStation, _ := h.repos.Station.GetByID(routeStations[len(routeStations)-1].StationID)
				if arrivalStation != nil {
					arrivalCity = arrivalStation.City
				}
				if routeStations[len(routeStations)-1].ArrivalTime != nil {
					arrivalTime = routeStations[len(routeStations)-1].ArrivalTime.Format("15:04")
				}
			}
		}
	}

	response := OrderResponse{
		ID:           order.ID,
		UserID:       order.UserID,
		RouteID:      order.RouteID,
		RouteName:    routeName,
		TrainNumber:  trainNumber,
		TrainType:    trainType,
		DepartureCity: departureCity,
		ArrivalCity:   arrivalCity,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		CreatedAt:    order.CreatedAt.Format(time.RFC3339),
		Status:       order.Status,
		TotalAmount:  order.TotalAmount,
		Tickets: []TicketResponse{
			{
				ID:           ticket.ID,
				TicketNumber: ticket.TicketNumber,
				SeatNumber:   &seat.Number,
				CarriageNumber: func() *int { if carriage != nil { n := carriage.Number; return &n }; return nil }(),
				Price:        ticket.Price,
			},
		},
	}

	c.JSON(http.StatusCreated, response)
}

// GetOrders gets user's orders
// @Summary Get user orders
// @Description Get list of orders for current user
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} OrderResponse
// @Router /orders [get]
func (h *Handlers) GetOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int64)

	orders, err := h.repos.Order.GetByUserID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	var responses []OrderResponse
	for _, order := range orders {
		tickets, _ := h.repos.Ticket.GetByOrderID(order.ID)
		var ticketResponses []TicketResponse
		for _, ticket := range tickets {
			var seatNumber, carriageNumber *int
			if ticket.SeatID != nil {
				seat, _ := h.repos.Seat.GetByID(*ticket.SeatID)
				if seat != nil {
					seatNumber = &seat.Number
					carriage, _ := h.repos.Carriage.GetByID(seat.CarriageID)
					if carriage != nil {
						n := carriage.Number
						carriageNumber = &n
					}
				}
			}
			ticketResponses = append(ticketResponses, TicketResponse{
				ID:             ticket.ID,
				TicketNumber:   ticket.TicketNumber,
				SeatNumber:     seatNumber,
				CarriageNumber: carriageNumber,
				Price:          ticket.Price,
			})
		}

		// Get route information
		var routeName, trainNumber, trainType, departureCity, arrivalCity, departureTime, arrivalTime string
		if order.RouteID != nil {
			route, _ := h.repos.Route.GetByID(*order.RouteID)
			if route != nil {
				routeName = route.Name
				train, _ := h.repos.Train.GetByID(route.TrainID)
				if train != nil {
					trainNumber = train.Number
					trainType = train.Type
				}
				routeStations, _ := h.repos.Route.GetStations(route.ID)
				if len(routeStations) > 0 {
					departureStation, _ := h.repos.Station.GetByID(routeStations[0].StationID)
					if departureStation != nil {
						departureCity = departureStation.City
					}
					if routeStations[0].DepartureTime != nil {
						departureTime = routeStations[0].DepartureTime.Format("15:04")
					}
					if len(routeStations) > 1 {
						arrivalStation, _ := h.repos.Station.GetByID(routeStations[len(routeStations)-1].StationID)
						if arrivalStation != nil {
							arrivalCity = arrivalStation.City
						}
						if routeStations[len(routeStations)-1].ArrivalTime != nil {
							arrivalTime = routeStations[len(routeStations)-1].ArrivalTime.Format("15:04")
						}
					}
				}
			}
		}

		responses = append(responses, OrderResponse{
			ID:           order.ID,
			UserID:       order.UserID,
			RouteID:      order.RouteID,
			RouteName:    routeName,
			TrainNumber:  trainNumber,
			TrainType:    trainType,
			DepartureCity: departureCity,
			ArrivalCity:   arrivalCity,
			DepartureTime: departureTime,
			ArrivalTime:   arrivalTime,
			CreatedAt:    order.CreatedAt.Format(time.RFC3339),
			Status:       order.Status,
			TotalAmount:  order.TotalAmount,
			Tickets:      ticketResponses,
		})
	}

	c.JSON(http.StatusOK, responses)
}

// DeleteOrder deletes an order
// @Summary Delete order
// @Description Delete an unpaid order
// @Tags Orders
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [delete]
func (h *Handlers) DeleteOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int64)

	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.repos.Order.GetByID(orderID)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Check ownership
	if order.UserID != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Only allow deletion of unpaid orders
	if order.Status != "PENDING" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only unpaid orders can be deleted"})
		return
	}

	if err := h.repos.Order.Delete(orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetOrder gets order details
// @Summary Get order details
// @Description Get detailed information about an order
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} OrderResponse
// @Router /orders/{id} [get]
func (h *Handlers) GetOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int64)

	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.repos.Order.GetByID(orderID)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Check ownership
	if order.UserID != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	tickets, _ := h.repos.Ticket.GetByOrderID(order.ID)
	var ticketResponses []TicketResponse
	for _, ticket := range tickets {
		var seatNumber, carriageNumber *int
		if ticket.SeatID != nil {
			seat, _ := h.repos.Seat.GetByID(*ticket.SeatID)
			if seat != nil {
				seatNumber = &seat.Number
				carriage, _ := h.repos.Carriage.GetByID(seat.CarriageID)
				if carriage != nil {
					n := carriage.Number
					carriageNumber = &n
				}
			}
		}
		ticketResponses = append(ticketResponses, TicketResponse{
			ID:             ticket.ID,
			TicketNumber:   ticket.TicketNumber,
			SeatNumber:     seatNumber,
			CarriageNumber: carriageNumber,
			Price:          ticket.Price,
		})
	}

	// Get route information
	var routeName, trainNumber, trainType, departureCity, arrivalCity, departureTime, arrivalTime string
	if order.RouteID != nil {
		route, _ := h.repos.Route.GetByID(*order.RouteID)
		if route != nil {
			routeName = route.Name
			train, _ := h.repos.Train.GetByID(route.TrainID)
			if train != nil {
				trainNumber = train.Number
				trainType = train.Type
			}
			routeStations, _ := h.repos.Route.GetStations(route.ID)
			if len(routeStations) > 0 {
				departureStation, _ := h.repos.Station.GetByID(routeStations[0].StationID)
				if departureStation != nil {
					departureCity = departureStation.City
				}
				if routeStations[0].DepartureTime != nil {
					departureTime = routeStations[0].DepartureTime.Format("15:04")
				}
				if len(routeStations) > 1 {
					arrivalStation, _ := h.repos.Station.GetByID(routeStations[len(routeStations)-1].StationID)
					if arrivalStation != nil {
						arrivalCity = arrivalStation.City
					}
					if routeStations[len(routeStations)-1].ArrivalTime != nil {
						arrivalTime = routeStations[len(routeStations)-1].ArrivalTime.Format("15:04")
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, OrderResponse{
		ID:           order.ID,
		UserID:       order.UserID,
		RouteID:      order.RouteID,
		RouteName:    routeName,
		TrainNumber:  trainNumber,
		TrainType:    trainType,
		DepartureCity: departureCity,
		ArrivalCity:   arrivalCity,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		CreatedAt:    order.CreatedAt.Format(time.RFC3339),
		Status:       order.Status,
		TotalAmount:  order.TotalAmount,
		Tickets:      ticketResponses,
	})
}

// PayOrder processes order payment
// @Summary Pay order
// @Description Process payment for an order
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body PaymentRequest true "Payment data"
// @Success 200 {object} PaymentResponse
// @Router /orders/{id}/pay [post]
func (h *Handlers) PayOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(int64)

	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.repos.Order.GetByID(orderID)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Check ownership
	if order.UserID != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update order status
	order.Status = "PAID"
	if err := h.repos.Order.Update(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	transactionID := fmt.Sprintf("TXN-%d-%d", orderID, time.Now().Unix())
	c.JSON(http.StatusOK, PaymentResponse{
		Status:        "PAID",
		TransactionID: transactionID,
	})
}

// CreateRoute creates a new route (Admin only)
// @Summary Create route
// @Description Create a new route (Admin only)
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateRouteRequest true "Route data"
// @Success 201 {object} models.Route
// @Router /admin/routes [post]
func (h *Handlers) CreateRoute(c *gin.Context) {
	var req CreateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route := &models.Route{
		Name:    req.Name,
		TrainID: req.TrainID,
	}
	if err := h.repos.Route.Create(route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create route"})
		return
	}

	c.JSON(http.StatusCreated, route)
}

// UpdateRoute updates a route (Admin only)
// @Summary Update route
// @Description Update a route (Admin only)
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Route ID"
// @Param request body UpdateRouteRequest true "Route data"
// @Success 200 {object} models.Route
// @Router /admin/routes/{id} [put]
func (h *Handlers) UpdateRoute(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID"})
		return
	}

	var req UpdateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.repos.Route.GetByID(id)
	if err != nil || route == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	if req.Name != "" {
		route.Name = req.Name
	}
	if req.TrainID != 0 {
		route.TrainID = req.TrainID
	}

	if err := h.repos.Route.Update(route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update route"})
		return
	}

	c.JSON(http.StatusOK, route)
}

// DeleteRoute deletes a route (Admin only)
// @Summary Delete route
// @Description Delete a route (Admin only)
// @Tags Admin
// @Security BearerAuth
// @Param id path int true "Route ID"
// @Success 204
// @Router /admin/routes/{id} [delete]
func (h *Handlers) DeleteRoute(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID"})
		return
	}

	if err := h.repos.Route.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete route"})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateTrain creates a new train (Admin only)
// @Summary Create train
// @Description Create a new train (Admin only)
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateTrainRequest true "Train data"
// @Success 201 {object} models.Train
// @Router /admin/trains [post]
func (h *Handlers) CreateTrain(c *gin.Context) {
	var req CreateTrainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	train := &models.Train{
		Number: req.Number,
		Type:   req.Type,
	}
	if err := h.repos.Train.Create(train); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create train"})
		return
	}

	c.JSON(http.StatusCreated, train)
}

// UpdateTrain updates a train (Admin only)
// @Summary Update train
// @Description Update a train (Admin only)
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Train ID"
// @Param request body UpdateTrainRequest true "Train data"
// @Success 200 {object} models.Train
// @Router /admin/trains/{id} [put]
func (h *Handlers) UpdateTrain(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid train ID"})
		return
	}

	var req UpdateTrainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	train, err := h.repos.Train.GetByID(id)
	if err != nil || train == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Train not found"})
		return
	}

	if req.Number != "" {
		train.Number = req.Number
	}
	if req.Type != "" {
		train.Type = req.Type
	}

	if err := h.repos.Train.Update(train); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update train"})
		return
	}

	c.JSON(http.StatusOK, train)
}

// DeleteTrain deletes a train (Admin only)
// @Summary Delete train
// @Description Delete a train (Admin only)
// @Tags Admin
// @Security BearerAuth
// @Param id path int true "Train ID"
// @Success 204
// @Router /admin/trains/{id} [delete]
func (h *Handlers) DeleteTrain(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid train ID"})
		return
	}

	if err := h.repos.Train.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete train"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllOrders gets all orders (Admin only)
// @Summary Get all orders
// @Description Get all orders in the system (Admin only)
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} OrderResponse
// @Router /admin/orders [get]
func (h *Handlers) GetAllOrders(c *gin.Context) {
	orders, err := h.repos.Order.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	var responses []OrderResponse
	for _, order := range orders {
		tickets, _ := h.repos.Ticket.GetByOrderID(order.ID)
		var ticketResponses []TicketResponse
		for _, ticket := range tickets {
			var seatNumber, carriageNumber *int
			if ticket.SeatID != nil {
				seat, _ := h.repos.Seat.GetByID(*ticket.SeatID)
				if seat != nil {
					seatNumber = &seat.Number
					carriage, _ := h.repos.Carriage.GetByID(seat.CarriageID)
					if carriage != nil {
						n := carriage.Number
						carriageNumber = &n
					}
				}
			}
			ticketResponses = append(ticketResponses, TicketResponse{
				ID:             ticket.ID,
				TicketNumber:   ticket.TicketNumber,
				SeatNumber:     seatNumber,
				CarriageNumber: carriageNumber,
				Price:          ticket.Price,
			})
		}

		// Get route information
		var routeName, trainNumber, trainType, departureCity, arrivalCity, departureTime, arrivalTime string
		if order.RouteID != nil {
			route, _ := h.repos.Route.GetByID(*order.RouteID)
			if route != nil {
				routeName = route.Name
				train, _ := h.repos.Train.GetByID(route.TrainID)
				if train != nil {
					trainNumber = train.Number
					trainType = train.Type
				}
				routeStations, _ := h.repos.Route.GetStations(route.ID)
				if len(routeStations) > 0 {
					departureStation, _ := h.repos.Station.GetByID(routeStations[0].StationID)
					if departureStation != nil {
						departureCity = departureStation.City
					}
					if routeStations[0].DepartureTime != nil {
						departureTime = routeStations[0].DepartureTime.Format("15:04")
					}
					if len(routeStations) > 1 {
						arrivalStation, _ := h.repos.Station.GetByID(routeStations[len(routeStations)-1].StationID)
						if arrivalStation != nil {
							arrivalCity = arrivalStation.City
						}
						if routeStations[len(routeStations)-1].ArrivalTime != nil {
							arrivalTime = routeStations[len(routeStations)-1].ArrivalTime.Format("15:04")
						}
					}
				}
			}
		}

		responses = append(responses, OrderResponse{
			ID:           order.ID,
			UserID:       order.UserID,
			RouteID:      order.RouteID,
			RouteName:    routeName,
			TrainNumber:  trainNumber,
			TrainType:    trainType,
			DepartureCity: departureCity,
			ArrivalCity:   arrivalCity,
			DepartureTime: departureTime,
			ArrivalTime:   arrivalTime,
			CreatedAt:    order.CreatedAt.Format(time.RFC3339),
			Status:       order.Status,
			TotalAmount:  order.TotalAmount,
			Tickets:      ticketResponses,
		})
	}

	c.JSON(http.StatusOK, responses)
}


