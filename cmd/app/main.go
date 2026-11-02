package main

import (
	"log"
	"net/http"

	"github.com/enolgor/lurra/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	c := controller.NewController()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/todos", func(r chi.Router) {
			r.Get("/", c.ListTodos)
			r.Post("/", c.AddTodo)
			r.Get("/{id}", c.ShowTodo)
			r.Patch("/{id}", c.UpdateTodo)
			r.Delete("/{id}", c.DeleteTodo)
		})
	})
	r.Handle("/*", http.FileServer(http.Dir("web/dist")))

	log.Fatal(http.ListenAndServe(":8080", r))
}
