package repository

import (
	"database/sql"
	"github.com/project13/backend-stealthisproject/internal/models"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	query := `INSERT INTO orders (user_id, route_id, status, total_amount) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRow(query, order.UserID, order.RouteID, order.Status, order.TotalAmount).Scan(&order.ID, &order.CreatedAt)
}

func (r *orderRepository) GetByID(id int64) (*models.Order, error) {
	order := &models.Order{}
	query := `SELECT id, user_id, route_id, created_at, status, total_amount FROM orders WHERE id = $1`
	var routeID sql.NullInt64
	err := r.db.QueryRow(query, id).Scan(&order.ID, &order.UserID, &routeID, &order.CreatedAt, &order.Status, &order.TotalAmount)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if routeID.Valid {
		order.RouteID = &routeID.Int64
	}
	return order, err
}

func (r *orderRepository) GetByUserID(userID int64) ([]models.Order, error) {
	query := `SELECT id, user_id, route_id, created_at, status, total_amount FROM orders WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var routeID sql.NullInt64
		if err := rows.Scan(&order.ID, &order.UserID, &routeID, &order.CreatedAt, &order.Status, &order.TotalAmount); err != nil {
			return nil, err
		}
		if routeID.Valid {
			order.RouteID = &routeID.Int64
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *orderRepository) GetAll() ([]models.Order, error) {
	query := `SELECT id, user_id, route_id, created_at, status, total_amount FROM orders ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var routeID sql.NullInt64
		if err := rows.Scan(&order.ID, &order.UserID, &routeID, &order.CreatedAt, &order.Status, &order.TotalAmount); err != nil {
			return nil, err
		}
		if routeID.Valid {
			order.RouteID = &routeID.Int64
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *orderRepository) Update(order *models.Order) error {
	query := `UPDATE orders SET status = $1, total_amount = $2 WHERE id = $3`
	_, err := r.db.Exec(query, order.Status, order.TotalAmount, order.ID)
	return err
}

func (r *orderRepository) Delete(id int64) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *orderRepository) DeleteExpiredPending(maxAgeMinutes int) error {
	query := `DELETE FROM orders WHERE status = 'PENDING' AND created_at < NOW() - INTERVAL '1 minute' * $1`
	_, err := r.db.Exec(query, maxAgeMinutes)
	return err
}

