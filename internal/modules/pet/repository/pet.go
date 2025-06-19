package repository

type Category struct {
	Id int `gorm:"primarykey"`
}

type Tag struct {
	Id int `gorm:"primarykey"`
}

type Pet struct {
	Id         int `gorm:"primarykey"`
	Name       string
	CategoryID int
	Category   Category `gorm:"foreignKey:CategoryID"`
	PhotoUrls  []string `gorm:"-"`
	Tag        []Tag    `gorm:"many2many:pet_tags"`
	Status     string
}
