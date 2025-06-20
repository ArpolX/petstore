package controller

type Order struct {
	Id       int    `json:"id" validate:"required,min=1"`
	PetId    int    `json:"petId" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
	ShipDate string `json:"shipDate" validate:"required"`
	Status   string `json:"status" validate:"required"`
	Complete bool   `json:"complete" validate:"required"`
}
