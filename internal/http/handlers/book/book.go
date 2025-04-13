package book

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/mangesh-shinde/booklib/internal/models"
	"github.com/mangesh-shinde/booklib/internal/storage"
	"github.com/mangesh-shinde/booklib/internal/utils/response"
)

var (
	BooksRe    = regexp.MustCompile("^/api/v1/books/*$")
	BookByIdRe = regexp.MustCompile("^/api/v1/books/([0-9]+)$")
)

type BookHandler struct {
	Storage storage.Storage
}

func (b *BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("url path", slog.String("url", r.URL.Path))
	switch {
	case r.Method == http.MethodPost:
		b.New(w, r)
		return
	case r.Method == http.MethodDelete:
		b.Delete(w, r)
		return
	case r.Method == http.MethodPut:
		b.Update(w, r)
		return
	case r.Method == http.MethodGet && BookByIdRe.MatchString(r.URL.Path):
		b.GetBookById(w, r)
		return
	case r.Method == http.MethodGet && BooksRe.MatchString(r.URL.Path):
		b.GetBooks(w, r)
		return
	default:
		return
	}
}

func (b *BookHandler) New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slog.Info("creating a book")
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if errors.Is(err, io.EOF) {
		response.SendError(w, http.StatusBadRequest, fmt.Errorf("empty body"))
		return
	}

	if err != nil {
		response.SendError(w, http.StatusBadRequest, err)
		return
	}

	// validate inputs before sending response
	if err := validator.New().Struct(book); err != nil {
		validateErrs := err.(validator.ValidationErrors)
		resp := response.ValidateErrors(validateErrs)
		response.WriteJsonResponse(w, http.StatusBadGateway, resp)
		return
	}

	bookId, err := b.Storage.CreateBook(book.Name, book.Author, book.PublicationDate, book.Price)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, fmt.Errorf("error while creating book"))
		return
	}

	response.WriteJsonResponse(w, http.StatusCreated, map[string]int64{"book_id": bookId})

}

func (b *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slog.Info("deleting a book")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book deleted"))
}

func (b *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	slog.Info("updating a book")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book updated"))
}

func (b *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	slog.Info("fetching books")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Books list here"))
}

func (b *BookHandler) GetBookById(w http.ResponseWriter, r *http.Request) {
	matches := BookByIdRe.FindStringSubmatch(r.URL.Path)
	slog.Info("fetching book", slog.String("id", matches[1]))
	bookId, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, err)
	}

	book, err := b.Storage.GetBook(bookId)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, err)
		return
	}

	response.WriteJsonResponse(w, http.StatusOK, book)
}
