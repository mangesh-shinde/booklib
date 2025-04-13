package sqlite

import (
	"database/sql"

	"github.com/mangesh-shinde/booklib/internal/config"
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
