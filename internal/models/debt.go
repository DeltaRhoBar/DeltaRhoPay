package models

type Debt struct {
	Resident_id int
	Resident Resident 		
	Moved string
}
