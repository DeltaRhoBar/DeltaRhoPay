package services

import (
	"database/sql"
	"deltapay/internal/models"
	"os"
	"path"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	GetResidents() ([]models.Resident, error)
	GetAllResidents() ([]models.Resident, error)
	GetBeverages() ([]models.Beverage, error)
	GetDebts() ([]models.Debt, error)
	GetOrders(int) ([]models.Order, error)
	AddResidentIfNotOccupied(int, int, string, string) (bool, error)
	AddResidentReplace(int, int, string, string) error
	UpdateResident(int, int, int, string, string) error
	AddBeverage(string, int) error
	RemoveBeverage(string) error
	AddOrder(string, int, int, int) error
	CheckOut() error
	SetMessage(string) error 
	GetMessage() (string, error)
	Pay(int) error
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
		r.id,
		r.r_floor,
		r.r_nr,
		r.name,
		r.telephone,
		COALESCE(r.removed_on, 'Still living here') AS moved_out_on,
		COALESCE(SUM(b.price * o.amount), 0) AS total_cost
		FROM residents r
		LEFT JOIN orders o
		ON o.resident_id = r.id
		AND o.paid_on IS NULL          
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
		if err := rows.Scan(&resident.Id, &resident.Room.Floor, &resident.Room.Nr, &resident.Name, &resident.Telephone, &resident.Moved, &resident.Debt); err != nil {
			return nil, err
		}
		residents = append(residents, resident)
	}
	return residents, nil
}

func (s *Sqlite) GetAllResidents() ([]models.Resident, error) {
	const query = `
		SELECT
		r.id,
		r.r_floor,
		r.r_nr,
		r.name,
		r.telephone,
		COALESCE(r.removed_on, 'Still living here') AS moved_out_on,
		COALESCE(SUM(b.price * o.amount), 0) AS total_cost
		FROM residents r
		LEFT JOIN orders o
		ON o.resident_id = r.id
		AND o.paid_on IS NULL          
		LEFT JOIN beverages b
		ON b.id = o.beverage_id
		GROUP BY r.id, r.r_floor, r.r_nr, r.name
		ORDER BY 
		r.removed_on IS NOT NULL,
		r.removed_on ASC,
		r.r_floor, r.r_nr;
		`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	residents := make([]models.Resident, 0)
	for rows.Next() {
		resident := models.Resident{}
		if err := rows.Scan(&resident.Id, &resident.Room.Floor, &resident.Room.Nr, &resident.Name, &resident.Telephone, &resident.Moved, &resident.Debt); err != nil {
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

func (s *Sqlite) GetDebts() ([]models.Debt, error) {
	const query = `
		SELECT
		r.id,
		r.name,
		r.r_floor,
		r.r_nr,
		COALESCE(r.removed_on, 'Still living here') AS moved_out_on,
		SUM(b.price * o.amount) AS unpaid_total
		FROM residents r
		JOIN orders o ON r.id = o.resident_id
		JOIN beverages b ON b.id = o.beverage_id
		WHERE o.paid_on IS NULL
		AND o.date < (SELECT MAX(date) FROM checkouts)
		GROUP BY r.id
		HAVING SUM(b.price * o.amount) > 0
		ORDER BY r.r_floor ASC, r.r_nr ASC;
		`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	debts := make([]models.Debt, 0)
	for rows.Next() {
		debt := models.Debt{}
		if err := rows.Scan(&debt.Resident.Id, &debt.Resident.Name, &debt.Resident.Room.Floor, &debt.Resident.Room.Nr, &debt.Resident.Moved, &debt.Resident.Debt ); err != nil {
			return nil, err
		}
		debts = append(debts, debt)
	}
	return debts, nil
}



func (s *Sqlite) GetOrders(page int) ([]models.Order, error) {
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
		ORDER BY o.date ASC
		LIMIT 50
		OFFSET ?;
		`
	offset := 50*(page-1)
	rows, err := s.db.Query(query, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		order := models.Order{}
		if err := rows.Scan(&order.Resident, &order.R_floor, &order.R_nr, &order.Beverage, &order.Amount, &order.Price, &order.Date, &order.Paid); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (s *Sqlite) AddResidentIfNotOccupied(r_floor int, r_nr int, name string, telephone string) (bool, error) {
	query := `INSERT INTO residents (r_floor, r_nr, name, telephone) VALUES (?, ?, ?, ?);`
	_, err := s.db.Exec(query, r_floor, r_nr, name, telephone)
	if err != nil {
		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (s *Sqlite) AddResidentReplace(r_floor int, r_nr int, name string, telephone string) error {
	query_remove := `
	UPDATE residents
	SET removed_on = CURRENT_TIMESTAMP
	WHERE r_floor = ?
	AND r_nr = ?
	AND removed_on IS NULL;`
	query_add := `
	INSERT INTO residents 
	(r_floor, r_nr, name, telephone) 
	VALUES (?, ?, ?, ?);`

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(query_remove, r_floor, r_nr, name)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(query_add, r_floor, r_nr, name, telephone)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit() 
}

func (s *Sqlite) UpdateResident(id int, r_floor int, r_nr int, name string, telephone string) error {
	query := `
	UPDATE residents SET
	name = ?,
	r_floor = ?,
	r_nr = ?,
	telephone = ?
	WHERE id = ?;`
	_, err := s.db.Exec(query, name, r_floor, r_nr, telephone, id)
	return err
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

func (s *Sqlite) CheckOut() error {
	query := `INSERT INTO checkouts DEFAULT VALUES;`
	_, err := s.db.Exec(query) 
	return err
}

func (s *Sqlite) SetMessage(message string) error {
	query := `INSERT OR REPLACE 
	INTO key_value (key, value)
	VALUES ('message', ?);`
	_, err := s.db.Exec(query, message) 
	return err;
}

func (s *Sqlite) GetMessage() (string, error) {
	query := `SELECT value 
	from key_value 
	WHERE key = 'message';`

	result := ""
	err := s.db.QueryRow(query).Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *Sqlite) Pay(id int) error {
	query := `UPDATE orders
	SET paid_on = CURRENT_TIMESTAMP
	WHERE resident_id = ?                        
	AND paid_on IS NULL                      
	AND date < (SELECT MAX(date) FROM checkouts);  `	
	_, err := s.db.Exec(query, id)
	return err
}
