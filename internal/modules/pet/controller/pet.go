package controller

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Pet struct {
	Id        int      `json:"id" validate:"required"`
	Name      string   `json:"name" validate:"required"`
	Category  Category `json:"category" validate:"required"`
	PhotoUrls []string `json:"photoUrls"`
	Tag       []Tag    `json:"tags" validate:"required"`
	Status    string   `json:"status" validate:"required"`
}

type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
