package repository

import (
	"database/sql"
	"github.com/project13/backend-stealthisproject/internal/models"
)

type carriageRepository struct {
	db *sql.DB
}

func NewCarriageRepository(db *sql.DB) CarriageRepository {
	return &carriageRepository{db: db}
}

func (r *carriageRepository) Create(carriage *models.Carriage) error {
	query := `INSERT INTO carriages (train_id, number, type) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(query, carriage.TrainID, carriage.Number, carriage.Type).Scan(&carriage.ID)
}

func (r *carriageRepository) GetByID(id int64) (*models.Carriage, error) {
	carriage := &models.Carriage{}
	query := `SELECT id, train_id, number, type FROM carriages WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&carriage.ID, &carriage.TrainID, &carriage.Number, &carriage.Type)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return carriage, err
}

func (r *carriageRepository) GetByTrainID(trainID int64) ([]models.Carriage, error) {
	query := `SELECT id, train_id, number, type FROM carriages WHERE train_id = $1 ORDER BY number`
	rows, err := r.db.Query(query, trainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carriages []models.Carriage
	for rows.Next() {
		var carriage models.Carriage
		if err := rows.Scan(&carriage.ID, &carriage.TrainID, &carriage.Number, &carriage.Type); err != nil {
			return nil, err
		}
		carriages = append(carriages, carriage)
	}
	return carriages, rows.Err()
}

