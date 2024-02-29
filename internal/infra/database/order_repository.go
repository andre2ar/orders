package database

import (
	"database/sql"

	"github.com/andre2ar/orders/internal/entity"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.DB.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
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
	rows, err := r.DB.Query("SELECT id, price, tax, final_price FROM orders")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var orders []entity.Order
	for rows.Next() {
		var currentOrder entity.Order
		err = rows.Scan(&currentOrder.ID, &currentOrder.Price, &currentOrder.Tax, &currentOrder.FinalPrice)
		if err != nil {
			return nil, err
		}

		orders = append(orders, currentOrder)
	}

	return orders, nil
}
