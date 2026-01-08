package repository

import (
	"database/sql"
	"github.com/project13/backend-stealthisproject/internal/models"
)

type trainRepository struct {
	db *sql.DB
}

func NewTrainRepository(db *sql.DB) TrainRepository {
	return &trainRepository{db: db}
}

func (r *trainRepository) Create(train *models.Train) error {
	query := `INSERT INTO trains (number, type) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, train.Number, train.Type).Scan(&train.ID)
}

func (r *trainRepository) GetByID(id int64) (*models.Train, error) {
	train := &models.Train{}
	query := `SELECT id, number, type FROM trains WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&train.ID, &train.Number, &train.Type)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return train, err
}

func (r *trainRepository) GetAll() ([]models.Train, error) {
	query := `SELECT id, number, type FROM trains ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trains []models.Train
	for rows.Next() {
		var train models.Train
		if err := rows.Scan(&train.ID, &train.Number, &train.Type); err != nil {
			return nil, err
		}
		trains = append(trains, train)
	}
	return trains, rows.Err()
}

func (r *trainRepository) Update(train *models.Train) error {
	query := `UPDATE trains SET number = $1, type = $2 WHERE id = $3`
	_, err := r.db.Exec(query, train.Number, train.Type, train.ID)
	return err
}

func (r *trainRepository) Delete(id int64) error {
	query := `DELETE FROM trains WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

