package storage

type Storage interface {
	CreateBook(bookName string, author string, publicationDate string, price float64) (int64, error)
}
