package models

type Book struct {
	Id              int64   `json:"id"`
	Name            string  `json:"book_name" validate:"required"`
	Author          string  `json:"author" validate:"required"`
	PublicationDate string  `json:"publication_date" validate:"required"`
	Price           float64 `json:"price" validate:"required,gt=0"`
}
