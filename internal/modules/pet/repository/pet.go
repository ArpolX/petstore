package repository

type Category struct {
	Id   int `gorm:"primarykey"`
	Name string
}

type Tag struct {
	Id   int `gorm:"primarykey"`
	Name string
}

type Pet struct {
	Id        int `gorm:"primarykey"`
	Name      string
	Category  Category `gorm:"goreignKey:CategotyId"`
	PhotoUrls []string `gorm:"-"`
	Tag       []Tag    `gorm:"many2many:pet_tags"`
	Status    string
}
