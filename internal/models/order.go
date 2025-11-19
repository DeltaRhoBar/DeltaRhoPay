package models

type Order struct {
	Resident string
	R_floor int
	R_nr int
	Beverage string
	Amount int
	Price int
	Date string
	Paid bool
}
