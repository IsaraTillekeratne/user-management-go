// @title User Management API
// @version 1.0
// @description REST API for User Management

package server

import (
	"net/http"

	_ "example.com/user-management/docs"
	"example.com/user-management/internal/db"
	"example.com/user-management/internal/handler"
	"example.com/user-management/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
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

	router.Get("/doc/*", httpSwagger.WrapHandler)

	return router
}
