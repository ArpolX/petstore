package service

type Order struct {
	Id       int
	PetId    int
	Quantity int
	ShipDate string
	Status   string
	Complete bool
}
