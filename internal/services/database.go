package services

import (
	"database/sql"
	"deltapay/internal/models"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mattn/go-sqlite3"
)

type Database interface {
	GetResidents() ([]models.Resident, error)
	GetBeverages() ([]models.Beverage, error)
	GetOrders() ([]models.Order, error)
	AddResidentIfNotOccupied(int, int, string) (bool, error)
	AddResidentReplace(int, int, string) error
	AddBeverage(string, int) error
	RemoveBeverage(string) error
	AddOrder(string, int, int, int) error
}

type Sqlite struct {
	db *sql.DB
}

func (s *Sqlite) executeSQLFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	query := string(bytes)
	_, err= s.db.Exec(query)
	return err
}

func NewSqlite() (*Sqlite, error) {
	db, err := sql.Open("sqlite3", path.Join("data", "sqlite.db"))
	if err != nil {
		return nil, err
	}

	s := &Sqlite{db: db}

	s.executeSQLFile(path.Join("data", "init.sql"))

	return s, nil
}

func (s *Sqlite) Close() {
	s.Close()
}

func (s *Sqlite) GetResidents() ([]models.Resident, error) {
	const query = `
		SELECT
		r.r_floor,
		r.r_nr,
		r.name,
		COALESCE(SUM(b.price * o.amount), 0) AS total_cost
		FROM residents r
		LEFT JOIN orders o
		ON o.resident_id = r.id
		LEFT JOIN beverages b
		ON b.id = o.beverage_id 
		WHERE r.removed_on IS NULL
		GROUP BY r.id, r.r_floor, r.r_nr, r.name
		ORDER BY r.r_floor, r.r_nr;
		`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	residents := make([]models.Resident, 0)
	for rows.Next() {
		resident := models.Resident{}
		if err := rows.Scan(&resident.Room.Floor, &resident.Room.Nr, &resident.Name, &resident.Debt); err != nil {
			return nil, err
		}
		residents = append(residents, resident)
	}
	return residents, nil
}

func (s *Sqlite) GetBeverages() ([]models.Beverage, error) {
	const query = `
		SELECT
		name,
		price
		FROM beverages
		WHERE removed_on IS NULL
		ORDER BY name ASC;
		`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beverages := make([]models.Beverage, 0)
	for rows.Next() {
		beverage := models.Beverage{}
		if err := rows.Scan(&beverage.Name, &beverage.Price); err != nil {
			return nil, err
		}
		beverages = append(beverages, beverage)
	}
	return beverages, nil
}

func (s *Sqlite) GetOrders() ([]models.Order, error) {
	const query = `
		SELECT
		r.name,
		r.r_floor,
		r.r_nr,
		b.name,
		o.amount,
		b.price,
		o.date,
		CASE WHEN o.paid_on IS NOT NULL THEN 1 ELSE 0 END AS paid
		FROM residents r
		JOIN orders o ON r.id = o.resident_id
		JOIN beverages b ON b.id = o.beverage_id
		ORDER BY o.date ASC;
		`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		order := models.Order{}
		if err := rows.Scan(&order.Resident, &order.R_Floor, &order.R_Nr, &order.Beverage, &order.Amount, &order.Price, &order.Date, &order.Paid); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return []models.Order{}, nil
}

func (s *Sqlite) AddResidentIfNotOccupied(r_floor int, r_nr int, name string) (bool, error) {
	query := `INSERT INTO residents (r_floor, r_nr, name) VALUES (?, ?, ?);`
	_, err := s.db.Exec(query, r_floor, r_nr, name)
	if err != nil {
		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (s *Sqlite) AddResidentReplace(r_floor int, r_nr int, name string) error {
	query_remove := `
	UPDATE residents
	SET removed_on = CURRENT_TIMESTAMP
	WHERE r_floor = ?
	AND r_nr = ?
	AND removed_on IS NULL;`
	query_add := `
	INSERT INTO residents 
	(r_floor, r_nr, name) 
	VALUES (?, ?, ?);`

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(query_remove, r_floor, r_nr, name)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(query_add, r_floor, r_nr, name)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit() 
}

func (s *Sqlite) AddBeverage(name string, price int) error {
	query := `INSERT INTO beverages (name, price) VALUES (?, ?);`	
	_, err := s.db.Exec(query, name, price)
	return err
}

func (s *Sqlite) RemoveBeverage(name string) error {
	query := `UPDATE beverages SET removed_on = CURRENT_TIMESTAMP WHERE name = ? AND removed_on IS NULL;`	
	_, err := s.db.Exec(query, name)
	return err
}

func (s *Sqlite) AddOrder(beverage_name string, amount int, r_floor int, r_nr int) error {
	query := `INSERT INTO orders (beverage_id, amount, resident_id)
	VALUES (
	(SELECT id FROM beverages
	WHERE name = ? 
	AND removed_on IS NULL),
	?,
	(SELECT id FROM residents
	WHERE r_floor = ?
	AND r_nr = ?
	AND removed_on IS NULL)
	);`
	_, err := s.db.Exec(query, beverage_name, amount, r_floor, r_nr)
	return err
}
