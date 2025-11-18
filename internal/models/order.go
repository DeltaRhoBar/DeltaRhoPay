package models

type Order struct {
	Resident string
	R_Floor int
	R_Nr int
	Beverage string
	Amount int
	Price int
	Date string
	Paid bool
}
