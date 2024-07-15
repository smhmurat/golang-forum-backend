package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang-forum-backend/controllers"
	"net/http"
)

func (a *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)

	if a.debug {
		mux.Use(middleware.Logger)
	}
	mux.Post("/auth/signup", controllers.SignUp)
	mux.Post("/auth/login", controllers.SignIn)

	return mux
}
