package book

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mangesh-shinde/booklib/internal/models"
	"github.com/mangesh-shinde/booklib/internal/utils/response"
)

type BookHandler struct{}

func (b *BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	case r.Method == http.MethodGet:
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
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Bad Request: Please validate input data")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&book)

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
