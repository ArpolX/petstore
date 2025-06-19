package controller

import "time"

type Order struct {
	Id       int       `json:"id" validate:"required"`
	PetId    int       `json:"petId" validate:"required"`
	Quantity int       `json:"quantity" validate:"required"`
	ShipDate time.Time `json:"shipDate" validate:"required"`
	Status   string    `json:"status" validate:"required"`
	Complete bool      `json:"complete" validate:"required"`
}
