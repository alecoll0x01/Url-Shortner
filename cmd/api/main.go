package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alecoll0x01/url-shortener/internal/handler"
	"github.com/alecoll0x01/url-shortener/internal/repository"
	"github.com/alecoll0x01/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	repo := repository.NewInMemoryRepository()
	svc := service.NewURLService(repo)
	h := handler.NewURLHandler(svc)

	r := chi.NewRouter()

	r.Use(middleware.Logger)    // loga cada request
	r.Use(middleware.Recoverer) // captura panics e retorna 500

	r.Post("/shorten", h.Shorten)
	r.Get("/{code}", h.Redirect)
	r.Get("/stats/{code}", h.Stats)
	r.Get("/urls", h.ListAll)

	port := ":8080"
	fmt.Printf("URL Shortener rodando em http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
