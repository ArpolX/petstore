package repository

import "time"

type Order struct {
	Id       int `gorm:"primarykey"`
	PetId    int
	Quantity int
	ShipDate time.Time
	Status   string
	Complete bool
}
