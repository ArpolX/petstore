package service

type Category struct {
	Id int
}

type Tag struct {
	Id int
}

type Pet struct {
	Id        int
	Name      string
	Category  Category
	PhotoUrls []string
	Tag       []Tag
	Status    string
}
