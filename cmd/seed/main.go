package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/project13/backend-stealthisproject/internal/database"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/railway_tickets?sslmode=disable"
	}

	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Seed trains
	seedTrains(db)

	// Seed carriages
	seedCarriages(db)

	// Seed seats
	seedSeats(db)

	// Seed stations
	seedStations(db)

	// Seed routes
	seedRoutes(db)

	// Seed route stations
	seedRouteStations(db)

	log.Println("Database seeded successfully!")
}

func seedTrains(db *sql.DB) {
	trains := []struct {
		number    string
		trainType string
	}{
		{"703Б", "Скоростной"},
		{"701Б", "Скоростной"},
		{"105Б", "Региональный"},
		{"107Б", "Региональный"},
	}

	for _, t := range trains {
		var id int64
		err := db.QueryRow("SELECT id FROM trains WHERE number = $1", t.number).Scan(&id)
		if err == sql.ErrNoRows {
			err = db.QueryRow("INSERT INTO trains (number, type) VALUES ($1, $2) RETURNING id", t.number, t.trainType).Scan(&id)
			if err != nil {
				log.Printf("Failed to insert train %s: %v", t.number, err)
			} else {
				log.Printf("Inserted train: %s (ID: %d)", t.number, id)
			}
		}
	}
}

func seedCarriages(db *sql.DB) {
	// Get train IDs
	var trainIDs []int64
	rows, err := db.Query("SELECT id FROM trains")
	if err != nil {
		log.Printf("Failed to get trains: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			trainIDs = append(trainIDs, id)
		}
	}

	// Create carriages for each train
	for _, trainID := range trainIDs {
		carriageTypes := []struct {
			number           int
			carriageType     string
			seatsPerCarriage int
		}{
			{1, "Плацкарт", 54},
			{2, "Купе", 36},
			{3, "СВ", 18},
		}

		for _, ct := range carriageTypes {
			var carriageID int64
			err := db.QueryRow("SELECT id FROM carriages WHERE train_id = $1 AND number = $2", trainID, ct.number).Scan(&carriageID)
			if err == sql.ErrNoRows {
				err = db.QueryRow("INSERT INTO carriages (train_id, number, type) VALUES ($1, $2, $3) RETURNING id", trainID, ct.number, ct.carriageType).Scan(&carriageID)
				if err != nil {
					log.Printf("Failed to insert carriage: %v", err)
					continue
				}

				// Create seats for this carriage
				for i := 1; i <= ct.seatsPerCarriage; i++ {
					_, err := db.Exec("INSERT INTO seats (carriage_id, number) VALUES ($1, $2)", carriageID, i)
					if err != nil {
						log.Printf("Failed to insert seat %d for carriage %d: %v", i, carriageID, err)
					}
				}
				log.Printf("Inserted carriage %d (type: %s) with %d seats for train %d", ct.number, ct.carriageType, ct.seatsPerCarriage, trainID)
			}
		}
	}
}

func seedSeats(db *sql.DB) {
	// Seats are created in seedCarriages
	log.Println("Seats created with carriages")
}

func seedStations(db *sql.DB) {
	stations := []struct {
		name string
		city string
	}{
		{"Минск-Пассажирский", "Минск"},
		{"Брест-Центральный", "Брест"},
		{"Гомель", "Гомель"},
		{"Витебск", "Витебск"},
		{"Гродно", "Гродно"},
		{"Могилев", "Могилев"},
	}

	for _, s := range stations {
		var id int64
		err := db.QueryRow("SELECT id FROM stations WHERE name = $1", s.name).Scan(&id)
		if err == sql.ErrNoRows {
			err = db.QueryRow("INSERT INTO stations (name, city) VALUES ($1, $2) RETURNING id", s.name, s.city).Scan(&id)
			if err != nil {
				log.Printf("Failed to insert station %s: %v", s.name, err)
			} else {
				log.Printf("Inserted station: %s (ID: %d)", s.name, id)
			}
		}
	}
}

func seedRoutes(db *sql.DB) {
	// Get train IDs
	var train703B, train701B, train105B, train107B int64
	db.QueryRow("SELECT id FROM trains WHERE number = '703Б' LIMIT 1").Scan(&train703B)
	db.QueryRow("SELECT id FROM trains WHERE number = '701Б' LIMIT 1").Scan(&train701B)
	db.QueryRow("SELECT id FROM trains WHERE number = '105Б' LIMIT 1").Scan(&train105B)
	db.QueryRow("SELECT id FROM trains WHERE number = '107Б' LIMIT 1").Scan(&train107B)

	routes := []struct {
		name    string
		trainID int64
		price   float64
	}{
		{"Минск - Брест", train703B, 28.00},
		{"Минск - Брест", train701B, 28.00},
		{"Минск - Гомель", train105B, 23.00},
		{"Минск - Витебск", train107B, 25.00},
	}

	for _, r := range routes {
		if r.trainID == 0 {
			continue // Skip if train not found
		}
		var id int64
		err := db.QueryRow("SELECT id FROM routes WHERE name = $1 AND train_id = $2", r.name, r.trainID).Scan(&id)
		if err == sql.ErrNoRows {
			err = db.QueryRow("INSERT INTO routes (name, train_id, price) VALUES ($1, $2, $3) RETURNING id", r.name, r.trainID, r.price).Scan(&id)
			if err != nil {
				log.Printf("Failed to insert route %s: %v", r.name, err)
			} else {
				log.Printf("Inserted route: %s (ID: %d, Price: %.2f BYN)", r.name, id, r.price)
			}
		} else {
			// Update price for existing routes
			_, err = db.Exec("UPDATE routes SET price = $1 WHERE id = $2", r.price, id)
			if err != nil {
				log.Printf("Failed to update price for route %s: %v", r.name, err)
			} else {
				log.Printf("Updated route price: %s (ID: %d, Price: %.2f BYN)", r.name, id, r.price)
			}
		}
	}
}

func seedRouteStations(db *sql.DB) {
	// Get station IDs
	var minskID, brestID, gomelID, vitebskID int64
	db.QueryRow("SELECT id FROM stations WHERE name = 'Минск-Пассажирский'").Scan(&minskID)
	db.QueryRow("SELECT id FROM stations WHERE name = 'Брест-Центральный'").Scan(&brestID)
	db.QueryRow("SELECT id FROM stations WHERE name = 'Гомель'").Scan(&gomelID)
	db.QueryRow("SELECT id FROM stations WHERE name = 'Витебск'").Scan(&vitebskID)

	// Get train IDs
	var train703B, train701B, train105B, train107B int64
	db.QueryRow("SELECT id FROM trains WHERE number = '703Б' LIMIT 1").Scan(&train703B)
	db.QueryRow("SELECT id FROM trains WHERE number = '701Б' LIMIT 1").Scan(&train701B)
	db.QueryRow("SELECT id FROM trains WHERE number = '105Б' LIMIT 1").Scan(&train105B)
	db.QueryRow("SELECT id FROM trains WHERE number = '107Б' LIMIT 1").Scan(&train107B)

	// Seed route stations for each route (by route name and train ID)
	routeStations := []struct {
		routeName     string
		trainID       int64
		departureTime string
		arrivalTime   string
		departureID   int64
		arrivalID     int64
	}{
		{"Минск - Брест", train703B, "08:00:00", "12:30:00", minskID, brestID},
		{"Минск - Брест", train701B, "14:00:00", "18:45:00", minskID, brestID},
		{"Минск - Гомель", train105B, "09:30:00", "14:15:00", minskID, gomelID},
		{"Минск - Витебск", train107B, "10:00:00", "15:30:00", minskID, vitebskID},
	}

	for _, rs := range routeStations {
		if rs.trainID == 0 || rs.departureID == 0 || rs.arrivalID == 0 {
			continue
		}

		// Get route ID by name and train ID
		var routeID int64
		err := db.QueryRow(`
			SELECT id FROM routes WHERE name = $1 AND train_id = $2
		`, rs.routeName, rs.trainID).Scan(&routeID)

		if err != nil {
			log.Printf("Failed to get route %s (train %d): %v", rs.routeName, rs.trainID, err)
			continue
		}

		// Check if route stations already exist
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM route_stations WHERE route_id = $1", routeID).Scan(&count)
		if err == nil && count == 0 {
			_, err = db.Exec(`
				INSERT INTO route_stations (route_id, station_id, departure_time, stop_order) 
				VALUES ($1, $2, $3, 1)
			`, routeID, rs.departureID, rs.departureTime)
			if err != nil {
				log.Printf("Failed to insert departure station for route %s (ID: %d): %v", rs.routeName, routeID, err)
				continue
			}

			_, err = db.Exec(`
				INSERT INTO route_stations (route_id, station_id, arrival_time, stop_order) 
				VALUES ($1, $2, $3, 2)
			`, routeID, rs.arrivalID, rs.arrivalTime)
			if err != nil {
				log.Printf("Failed to insert arrival station for route %s (ID: %d): %v", rs.routeName, routeID, err)
			} else {
				log.Printf("Inserted route stations for route %s (ID: %d)", rs.routeName, routeID)
			}
		} else {
			log.Printf("Route stations already exist for route %s (ID: %d)", rs.routeName, routeID)
		}
	}
}
