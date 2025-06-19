package controller

type Category struct {
	Id int `json:"id" validate:"required,min=1,max=2"`
}

type Tag struct {
	Id int `json:"id" validate:"required,min=1,max=3"`
}

type Pet struct {
	Id        int      `json:"id" validate:"required"`
	Name      string   `json:"name" validate:"required"`
	Category  Category `json:"category" validate:"required"`
	PhotoUrls []string `json:"photoUrls"`
	Tag       []Tag    `json:"tags" validate:"required,dive"`
	Status    string   `json:"status" validate:"required"`
}

type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
