package repository

import (
	"database/sql"
	"github.com/project13/backend-stealthisproject/internal/models"
)

type ticketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ticket *models.Ticket) error {
	query := `INSERT INTO tickets (order_id, seat_id, passenger_id, departure_date, price, ticket_number, status) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRow(query, ticket.OrderID, ticket.SeatID, ticket.PassengerID, ticket.DepartureDate, 
		ticket.Price, ticket.TicketNumber, ticket.Status).Scan(&ticket.ID)
}

func (r *ticketRepository) GetByID(id int64) (*models.Ticket, error) {
	ticket := &models.Ticket{}
	query := `SELECT id, order_id, seat_id, passenger_id, departure_date, price, ticket_number, status 
	          FROM tickets WHERE id = $1`
	var seatID, passengerID sql.NullInt64
	err := r.db.QueryRow(query, id).Scan(&ticket.ID, &ticket.OrderID, &seatID, &passengerID, 
		&ticket.DepartureDate, &ticket.Price, &ticket.TicketNumber, &ticket.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if seatID.Valid {
		ticket.SeatID = &seatID.Int64
	}
	if passengerID.Valid {
		ticket.PassengerID = &passengerID.Int64
	}
	return ticket, err
}

func (r *ticketRepository) GetByOrderID(orderID int64) ([]models.Ticket, error) {
	query := `SELECT id, order_id, seat_id, passenger_id, departure_date, price, ticket_number, status 
	          FROM tickets WHERE order_id = $1 ORDER BY id`
	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var ticket models.Ticket
		var seatID, passengerID sql.NullInt64
		if err := rows.Scan(&ticket.ID, &ticket.OrderID, &seatID, &passengerID, 
			&ticket.DepartureDate, &ticket.Price, &ticket.TicketNumber, &ticket.Status); err != nil {
			return nil, err
		}
		if seatID.Valid {
			ticket.SeatID = &seatID.Int64
		}
		if passengerID.Valid {
			ticket.PassengerID = &passengerID.Int64
		}
		tickets = append(tickets, ticket)
	}
	return tickets, rows.Err()
}

func (r *ticketRepository) Update(ticket *models.Ticket) error {
	query := `UPDATE tickets SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, ticket.Status, ticket.ID)
	return err
}

