package server

import (
	"net/http"

	"example.com/user-management/internal/db"
	"example.com/user-management/internal/handler"
	"example.com/user-management/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	dbConn := db.NewPostgres()
	userStore := store.NewUserStore(dbConn)
	userHandler := handler.NewUserHandler(userStore)

	router.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.GetAllUsers)
		r.Get("/{id}", userHandler.GetUserById)
		r.Patch("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})

	return router
}
