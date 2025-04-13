package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mangesh-shinde/booklib/internal/config"
	"github.com/mangesh-shinde/booklib/internal/http/handlers/book"
	"github.com/mangesh-shinde/booklib/internal/storage/sqlite"
)

func main() {

	// load config
	cfg := config.MustLoad()

	// setup database
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info(">> storage initialized")

	// setup router
	mux := http.NewServeMux()

	// book routes
	mux.Handle("/api/v1/books", &book.BookHandler{Storage: storage})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to booklib API"))
	})

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		// fmt.Printf("server is listening on address: %s\n", cfg.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	slog.Info("Server started on ", slog.String("address", cfg.Addr))
	<-done

	slog.Info("Shutting down server...")
	// graceful server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}
