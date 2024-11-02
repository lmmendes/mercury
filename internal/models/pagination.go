package models

type PaginationQuery struct {
	Limit  int `query:"limit" validate:"min=1,max=100"`
	Offset int `query:"offset" validate:"min=0"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		Total  int `json:"total"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"pagination"`
}
