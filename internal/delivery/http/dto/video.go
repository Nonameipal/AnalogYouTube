package dto

type CreateVideoRequest struct {
	Title        string
	Description  string
	CategoryName *string
}

type UpdateVideoRequest struct {
	Title        string
	Description  string
	CategoryName *string
}
