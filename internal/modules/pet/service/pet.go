package service

type Category struct {
	Id   int
	Name string
}

type Tag struct {
	Id   int
	Name string
}

type Pet struct {
	Id        int
	Name      string
	Category  Category
	PhotoUrls []string
	Tag       []Tag
	Status    string
}

type ApiResponse struct {
	Code    int
	Type    string
	Message string
}
