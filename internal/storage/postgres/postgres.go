package postgres

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/sema0205/WB_L0/internal/models"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storageAuth string) (*Storage, error) {

	db, err := sql.Open("postgres", storageAuth)
	if err != nil {
		return nil, err
	}

	createDbRequest, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS ORDERS(
	    orderId VARCHAR PRIMARY KEY,
	    orderInfo JSONB);
	`)

	if err != nil {
		return nil, err
	}

	test, err := createDbRequest.Exec()
	if err != nil {
		return nil, err
	}
	test.LastInsertId()

	return &Storage{db: db}, nil
}

func (s *Storage) SaveOrder(orderInfo models.OrderInfo) error {

	insertOrderInfo, err := s.db.Prepare("INSERT INTO orders(orderId, orderInfo) VALUES($1, $2)")
	if err != nil {
		return err
	}

	jsonString, _ := json.Marshal(orderInfo)

	_, err = insertOrderInfo.Exec(orderInfo.OrderUID, jsonString)

	return nil
}

func (s *Storage) GetOrders() ([]models.OrderInfo, error) {

	getOrderInfo, err := s.db.Prepare("SELECT * FROM orders")
	if err != nil {
		return []models.OrderInfo{}, err
	}

	var orderInfo []models.OrderInfo

	rows, err := getOrderInfo.Query()

	for rows.Next() {

		var uuid string
		var modelBytes []byte
		var model models.OrderInfo

		err = rows.Scan(&uuid, &modelBytes)
		json.Unmarshal(modelBytes, &model)
		orderInfo = append(orderInfo, model)
	}

	if err != nil {
		return []models.OrderInfo{}, err
	}

	return orderInfo, nil
}

func (s *Storage) RestoreCache() map[string]models.OrderInfo {

	cache := make(map[string]models.OrderInfo)

	orderInfo, _ := s.GetOrders()

	for _, info := range orderInfo {
		cache[info.OrderUID] = info
	}

	return cache
}
