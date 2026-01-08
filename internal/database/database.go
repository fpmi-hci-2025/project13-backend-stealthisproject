package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createUsersTable,
		createPassengersTable,
		createTrainsTable,
		createCarriagesTable,
		createSeatsTable,
		createStationsTable,
		createRoutesTable,
		createRouteStationsTable,
		createOrdersTable,
		createTicketsTable,
		createIndexes,
		// Add route_id column to orders table if it doesn't exist
		`ALTER TABLE orders ADD COLUMN IF NOT EXISTS route_id BIGINT REFERENCES routes(id) ON DELETE SET NULL`,
		// Add price column to routes table if it doesn't exist
		`ALTER TABLE routes ADD COLUMN IF NOT EXISTS price DECIMAL(10, 2) NOT NULL DEFAULT 20.00`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'PASSENGER',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
`

const createPassengersTable = `
CREATE TABLE IF NOT EXISTS passengers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    passport_data VARCHAR(50)
);
`

const createTrainsTable = `
CREATE TABLE IF NOT EXISTS trains (
    id BIGSERIAL PRIMARY KEY,
    number VARCHAR(20) NOT NULL UNIQUE,
    type VARCHAR(50)
);
`

const createCarriagesTable = `
CREATE TABLE IF NOT EXISTS carriages (
    id BIGSERIAL PRIMARY KEY,
    train_id BIGINT REFERENCES trains(id) ON DELETE CASCADE,
    number INTEGER NOT NULL,
    type VARCHAR(50)
);
`

const createSeatsTable = `
CREATE TABLE IF NOT EXISTS seats (
    id BIGSERIAL PRIMARY KEY,
    carriage_id BIGINT REFERENCES carriages(id) ON DELETE CASCADE,
    number INTEGER NOT NULL
);
`

const createStationsTable = `
CREATE TABLE IF NOT EXISTS stations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL
);
`

const createRoutesTable = `
CREATE TABLE IF NOT EXISTS routes (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    train_id BIGINT REFERENCES trains(id) ON DELETE SET NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 20.00
);
`

const createRouteStationsTable = `
CREATE TABLE IF NOT EXISTS route_stations (
    route_id BIGINT REFERENCES routes(id) ON DELETE CASCADE,
    station_id BIGINT REFERENCES stations(id) ON DELETE CASCADE,
    arrival_time TIME,
    departure_time TIME,
    stop_order INTEGER NOT NULL,
    PRIMARY KEY (route_id, station_id)
);
`

const createOrdersTable = `
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    route_id BIGINT REFERENCES routes(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00
);
`

const createTicketsTable = `
CREATE TABLE IF NOT EXISTS tickets (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES orders(id) ON DELETE CASCADE,
    seat_id BIGINT REFERENCES seats(id),
    passenger_id BIGINT REFERENCES passengers(id) ON DELETE SET NULL,
    departure_date DATE NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    ticket_number VARCHAR(50) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE'
);
`

const createIndexes = `
CREATE INDEX IF NOT EXISTS idx_tickets_order_id ON tickets(order_id);
CREATE INDEX IF NOT EXISTS idx_tickets_passenger_id ON tickets(passenger_id);
CREATE INDEX IF NOT EXISTS idx_route_stations_route_id ON route_stations(route_id);
`

