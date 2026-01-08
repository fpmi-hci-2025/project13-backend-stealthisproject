package repository

import (
	"database/sql"
	"github.com/project13/backend-stealthisproject/internal/models"
)

type Repositories struct {
	User      UserRepository
	Passenger PassengerRepository
	Train     TrainRepository
	Carriage  CarriageRepository
	Seat      SeatRepository
	Station   StationRepository
	Route     RouteRepository
	Order     OrderRepository
	Ticket    TicketRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:      NewUserRepository(db),
		Passenger: NewPassengerRepository(db),
		Train:     NewTrainRepository(db),
		Carriage:  NewCarriageRepository(db),
		Seat:      NewSeatRepository(db),
		Station:   NewStationRepository(db),
		Route:     NewRouteRepository(db),
		Order:     NewOrderRepository(db),
		Ticket:    NewTicketRepository(db),
	}
}

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
}

type PassengerRepository interface {
	Create(passenger *models.Passenger) error
	GetByID(id int64) (*models.Passenger, error)
	GetByUserID(userID int64) (*models.Passenger, error)
	Update(passenger *models.Passenger) error
}

type TrainRepository interface {
	Create(train *models.Train) error
	GetByID(id int64) (*models.Train, error)
	GetAll() ([]models.Train, error)
	Update(train *models.Train) error
	Delete(id int64) error
}

type CarriageRepository interface {
	Create(carriage *models.Carriage) error
	GetByID(id int64) (*models.Carriage, error)
	GetByTrainID(trainID int64) ([]models.Carriage, error)
}

type SeatRepository interface {
	Create(seat *models.Seat) error
	GetByID(id int64) (*models.Seat, error)
	GetByCarriageID(carriageID int64) ([]models.Seat, error)
	IsAvailable(seatID int64, date string) (bool, error)
}

type StationRepository interface {
	Create(station *models.Station) error
	GetByID(id int64) (*models.Station, error)
	GetByCity(city string) ([]models.Station, error)
	GetAll() ([]models.Station, error)
}

type RouteRepository interface {
	Create(route *models.Route) error
	GetByID(id int64) (*models.Route, error)
	Search(fromCity, toCity, date string) ([]models.Route, error)
	Update(route *models.Route) error
	Delete(id int64) error
	AddStation(routeID, stationID int64, arrivalTime, departureTime string, stopOrder int) error
	GetStations(routeID int64) ([]models.RouteStation, error)
}

type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id int64) (*models.Order, error)
	GetByUserID(userID int64) ([]models.Order, error)
	GetAll() ([]models.Order, error)
	Update(order *models.Order) error
	Delete(id int64) error
	DeleteExpiredPending(maxAgeMinutes int) error
}

type TicketRepository interface {
	Create(ticket *models.Ticket) error
	GetByID(id int64) (*models.Ticket, error)
	GetByOrderID(orderID int64) ([]models.Ticket, error)
	Update(ticket *models.Ticket) error
}

