package storage

import "github.com/mangesh-shinde/booklib/internal/models"

type Storage interface {
	CreateBook(bookName string, author string, publicationDate string, price float64) (int64, error)
	GetBook(id int64) (models.Book, error)
	DeleteBook(id int64) (int64, error)
	GetBooks() ([]models.Book, error)
}
