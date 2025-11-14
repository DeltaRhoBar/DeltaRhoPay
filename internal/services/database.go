package services

import (
	"database/sql"
	"deltapay/internal/models"
	"log"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	GetResidents() ([]models.Resident, error)
	AddResident(int, int, string) error
	CheckOccupation(int, int) (bool, error)
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
		r_floor,
		r_nr,
		name,
		IFNULL(d.amount, 0) AS debt
		FROM residents
		LEFT JOIN debts d ON residents.id = d.resident_id AND (d.date IS NULL OR d.date = '') ORDER BY r_floor ASC, r_nr ASC
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

func (s *Sqlite) AddResident(r_floor int, r_nr int, name string) error {
	query := `INSERT INTO residents (r_floor, r_nr, name) VALUES (?, ?, ?);`
	_, err := s.db.Exec(query, r_floor, r_nr, name)
	return err
}

func (s *Sqlite) CheckOccupation(r_floor int, r_nr int) (bool, error) {
	query := `SELECT EXISTS (
	SELECT 1
	FROM residents
	WHERE r_floor = ?
	AND r_nr = ?
	AND removed_on IS NULL
	) AS is_occupied;`

	row := s.db.QueryRow(query, r_floor, r_nr)
	var is_occupied bool
	err := row.Scan(&is_occupied)
	if err != nil {
		return false, err
	}
	log.Print(is_occupied)
	return is_occupied, nil
}
