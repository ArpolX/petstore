package service

import "time"

type Order struct {
	Id       int
	PetId    int
	Quantity int
	ShipDate time.Time
	Status   string
	Complete bool
}
