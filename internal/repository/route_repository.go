package repository

import (
	"database/sql"
	"time"
	"github.com/project13/backend-stealthisproject/internal/models"
)

type routeRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) RouteRepository {
	return &routeRepository{db: db}
}

func (r *routeRepository) Create(route *models.Route) error {
	query := `INSERT INTO routes (name, train_id, price) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(query, route.Name, route.TrainID, route.Price).Scan(&route.ID)
}

func (r *routeRepository) GetByID(id int64) (*models.Route, error) {
	route := &models.Route{}
	query := `SELECT id, name, train_id, price FROM routes WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&route.ID, &route.Name, &route.TrainID, &route.Price)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return route, err
}

func (r *routeRepository) Search(fromCity, toCity, date string) ([]models.Route, error) {
	query := `
		SELECT DISTINCT r.id, r.name, r.train_id, r.price
		FROM routes r
		INNER JOIN route_stations rs1 ON r.id = rs1.route_id
		INNER JOIN stations s1 ON rs1.station_id = s1.id
		INNER JOIN route_stations rs2 ON r.id = rs2.route_id
		INNER JOIN stations s2 ON rs2.station_id = s2.id
		WHERE s1.city = $1 AND s2.city = $2 AND rs1.stop_order < rs2.stop_order
	`
	rows, err := r.db.Query(query, fromCity, toCity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []models.Route
	for rows.Next() {
		var route models.Route
		if err := rows.Scan(&route.ID, &route.Name, &route.TrainID, &route.Price); err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	return routes, rows.Err()
}

func (r *routeRepository) Update(route *models.Route) error {
	query := `UPDATE routes SET name = $1, train_id = $2, price = $3 WHERE id = $4`
	_, err := r.db.Exec(query, route.Name, route.TrainID, route.Price, route.ID)
	return err
}

func (r *routeRepository) Delete(id int64) error {
	query := `DELETE FROM routes WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *routeRepository) AddStation(routeID, stationID int64, arrivalTime, departureTime string, stopOrder int) error {
	var arrTime, depTime interface{}
	if arrivalTime != "" {
		t, _ := time.Parse("15:04:05", arrivalTime)
		arrTime = t
	}
	if departureTime != "" {
		t, _ := time.Parse("15:04:05", departureTime)
		depTime = t
	}
	
	query := `INSERT INTO route_stations (route_id, station_id, arrival_time, departure_time, stop_order) 
	          VALUES ($1, $2, $3, $4, $5) 
	          ON CONFLICT (route_id, station_id) DO UPDATE SET 
	          arrival_time = EXCLUDED.arrival_time, 
	          departure_time = EXCLUDED.departure_time, 
	          stop_order = EXCLUDED.stop_order`
	_, err := r.db.Exec(query, routeID, stationID, arrTime, depTime, stopOrder)
	return err
}

func (r *routeRepository) GetStations(routeID int64) ([]models.RouteStation, error) {
	query := `
		SELECT route_id, station_id, arrival_time, departure_time, stop_order
		FROM route_stations
		WHERE route_id = $1
		ORDER BY stop_order
	`
	rows, err := r.db.Query(query, routeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routeStations []models.RouteStation
	for rows.Next() {
		var rs models.RouteStation
		var arrTime, depTime sql.NullTime
		if err := rows.Scan(&rs.RouteID, &rs.StationID, &arrTime, &depTime, &rs.StopOrder); err != nil {
			return nil, err
		}
		if arrTime.Valid {
			rs.ArrivalTime = &arrTime.Time
		}
		if depTime.Valid {
			rs.DepartureTime = &depTime.Time
		}
		routeStations = append(routeStations, rs)
	}
	return routeStations, rows.Err()
}

