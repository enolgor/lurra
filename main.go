package main

import (
	"log"
	"net/http"

	"github.com/enolgor/lurra/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:generate go tool -modfile=tools.mod swag init --v3.1 --outputTypes yaml
//go:generate npm --prefix web exec -- openapi-typescript docs/swagger.yaml -o web/src/schema.ts

//	@title			Todo API
//	@version		1.0
//	@description	A very simple todo API.

//	@servers.url	http://localhost:8080/api/v1

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
