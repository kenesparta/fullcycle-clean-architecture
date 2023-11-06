package database

import (
	"database/sql"

	"github.com/kenesparta/fullcycle-clean-architecture/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) List() ([]entity.Order, error) {
	rows, err := r.Db.Query("SELECT id, price, tax, final_price from orders")
	if err != nil {
		return nil, err
	}

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		scanErr := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if scanErr != nil {
			return nil, scanErr
		}
		orders = append(orders, order)
	}

	return orders, nil
}
