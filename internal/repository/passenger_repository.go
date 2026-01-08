package repository

import (
	"database/sql"
	"github.com/project13/backend-stealthisproject/internal/models"
)

type passengerRepository struct {
	db *sql.DB
}

func NewPassengerRepository(db *sql.DB) PassengerRepository {
	return &passengerRepository{db: db}
}

func (r *passengerRepository) Create(passenger *models.Passenger) error {
	query := `INSERT INTO passengers (user_id, first_name, last_name, passport_data) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(query, passenger.UserID, passenger.FirstName, passenger.LastName, passenger.PassportData).Scan(&passenger.ID)
}

func (r *passengerRepository) GetByID(id int64) (*models.Passenger, error) {
	passenger := &models.Passenger{}
	query := `SELECT id, user_id, first_name, last_name, passport_data FROM passengers WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&passenger.ID, &passenger.UserID, &passenger.FirstName, &passenger.LastName, &passenger.PassportData)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return passenger, err
}

func (r *passengerRepository) GetByUserID(userID int64) (*models.Passenger, error) {
	passenger := &models.Passenger{}
	query := `SELECT id, user_id, first_name, last_name, passport_data FROM passengers WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&passenger.ID, &passenger.UserID, &passenger.FirstName, &passenger.LastName, &passenger.PassportData)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return passenger, err
}

func (r *passengerRepository) Update(passenger *models.Passenger) error {
	query := `UPDATE passengers SET first_name = $1, last_name = $2, passport_data = $3 WHERE id = $4`
	_, err := r.db.Exec(query, passenger.FirstName, passenger.LastName, passenger.PassportData, passenger.ID)
	return err
}

