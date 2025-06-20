package repository

type Order struct {
	Id       int `gorm:"primarykey"`
	PetId    int
	Quantity int
	ShipDate string
	Status   string
	Complete bool
}

type Pet struct {
	Id         int `gorm:"primarykey"`
	Name       string
	CategoryID int
	Category   Category `gorm:"foreignKey:CategoryID"`
	PhotoUrl   string
	Tag        []Tag `gorm:"many2many:pet_tags"`
	Status     string
}

type Category struct {
	Id int `gorm:"primarykey"`
}

type Tag struct {
	Id int `gorm:"primarykey"`
}
