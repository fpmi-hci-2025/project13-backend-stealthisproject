package repository

import (
	"database/sql"
	"github.com/project13/backend-stealthisproject/internal/models"
)

type stationRepository struct {
	db *sql.DB
}

func NewStationRepository(db *sql.DB) StationRepository {
	return &stationRepository{db: db}
}

func (r *stationRepository) Create(station *models.Station) error {
	query := `INSERT INTO stations (name, city) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, station.Name, station.City).Scan(&station.ID)
}

func (r *stationRepository) GetByID(id int64) (*models.Station, error) {
	station := &models.Station{}
	query := `SELECT id, name, city FROM stations WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&station.ID, &station.Name, &station.City)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return station, err
}

func (r *stationRepository) GetByCity(city string) ([]models.Station, error) {
	query := `SELECT id, name, city FROM stations WHERE city = $1 ORDER BY name`
	rows, err := r.db.Query(query, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stations []models.Station
	for rows.Next() {
		var station models.Station
		if err := rows.Scan(&station.ID, &station.Name, &station.City); err != nil {
			return nil, err
		}
		stations = append(stations, station)
	}
	return stations, rows.Err()
}

func (r *stationRepository) GetAll() ([]models.Station, error) {
	query := `SELECT id, name, city FROM stations ORDER BY city, name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stations []models.Station
	for rows.Next() {
		var station models.Station
		if err := rows.Scan(&station.ID, &station.Name, &station.City); err != nil {
			return nil, err
		}
		stations = append(stations, station)
	}
	return stations, rows.Err()
}

