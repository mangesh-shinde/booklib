package models

type Book struct {
	Id              int64   `json:"id"`
	Name            string  `json:"book_name"`
	Author          string  `json:"author"`
	PublicationDate string  `json:"publication_date"`
	Price           float64 `json:"price"`
}
