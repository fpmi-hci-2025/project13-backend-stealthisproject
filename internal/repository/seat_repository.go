package repository

import (
	"database/sql"
	"github.com/project13/backend-stealthisproject/internal/models"
)

type seatRepository struct {
	db *sql.DB
}

func NewSeatRepository(db *sql.DB) SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) Create(seat *models.Seat) error {
	query := `INSERT INTO seats (carriage_id, number) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, seat.CarriageID, seat.Number).Scan(&seat.ID)
}

func (r *seatRepository) GetByID(id int64) (*models.Seat, error) {
	seat := &models.Seat{}
	query := `SELECT id, carriage_id, number FROM seats WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&seat.ID, &seat.CarriageID, &seat.Number)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return seat, err
}

func (r *seatRepository) GetByCarriageID(carriageID int64) ([]models.Seat, error) {
	query := `SELECT id, carriage_id, number FROM seats WHERE carriage_id = $1 ORDER BY number`
	rows, err := r.db.Query(query, carriageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var seat models.Seat
		if err := rows.Scan(&seat.ID, &seat.CarriageID, &seat.Number); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	return seats, rows.Err()
}

func (r *seatRepository) IsAvailable(seatID int64, date string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM tickets 
		WHERE seat_id = $1 AND departure_date = $2 AND status = 'ACTIVE'
	`
	var count int
	err := r.db.QueryRow(query, seatID, date).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

