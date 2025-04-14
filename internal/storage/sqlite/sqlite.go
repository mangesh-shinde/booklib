package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/mangesh-shinde/booklib/internal/config"
	"github.com/mangesh-shinde/booklib/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	// this function will open a db connection and return an instance of Sqlite
	db, err := sql.Open("sqlite3", cfg.Storage)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		book_name TEXT,
		author TEXT,
		price REAL,
		publication_date TEXT
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateBook(bookName string, author string, publicationDate string, price float64) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO books (book_name, author, price, publication_date) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(bookName, author, price, publicationDate)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil

}

func (s *Sqlite) GetBook(id int64) (models.Book, error) {
	stmt, err := s.Db.Prepare("SELECT id, book_name, author, price, publication_date FROM books where id=? LIMIT 1")
	if err != nil {
		return models.Book{}, err
	}
	defer stmt.Close()

	var book models.Book
	err = stmt.QueryRow(id).Scan(&book.Id, &book.Name, &book.Author, &book.Price, &book.PublicationDate)
	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func (s *Sqlite) GetBooks() ([]models.Book, error) {
	stmt, err := s.Db.Prepare("SELECT id, book_name, author, price, publication_date FROM books")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.Id, &book.Name, &book.Author, &book.Price, &book.PublicationDate)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (s *Sqlite) DeleteBook(id int64) (int64, error) {
	stmt, err := s.Db.Prepare("DELETE FROM books WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsDeleted, nil
}

func (s *Sqlite) UpdateBook(id int64, bookName string, author string, publicationDate string, price float64) (int64, error) {
	// check if book for given id exists in database
	stmt, err := s.Db.Prepare("SELECT id FROM books WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var bookId int64
	err = stmt.QueryRow(id).Scan(&bookId)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("book not found with id=%s", fmt.Sprint(id))
	}

	if err != nil {
		return 0, err
	}

	// update book details
	updateStmt, err := s.Db.Prepare("UPDATE books SET book_name=?, author=?, price=?, publication_date=? WHERE id=?")
	if err != nil {
		return 0, err
	}
	defer updateStmt.Close()

	result, err := updateStmt.Exec(bookName, author, price, publicationDate, bookId)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
