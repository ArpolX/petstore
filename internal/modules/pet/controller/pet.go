package controller

type Category struct {
	Id int `json:"id" validate:"required,min=1,max=2"`
}

type Tag struct {
	Id int `json:"id" validate:"required,min=1,max=3"`
}

type Pet struct {
	Id       int      `json:"id" validate:"required"`
	Name     string   `json:"name" validate:"required"`
	Category Category `json:"category" validate:"required"`
	Tag      []Tag    `json:"tags" validate:"required,dive"`
	Status   string   `json:"status" validate:"required"`
}

type OutputPet struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Category_id int         `json:"category_id"`
	PhotoUrl    string      `json:"photoUrl"`
	Status      string      `json:"status"`
	OutputTag   []OutputTag `json:"tags"`
}

type OutputTag struct {
	Id int `json:"id"`
}

type OutputPetArray struct {
	OutputPet []OutputPet `json:"pets"`
}

type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
