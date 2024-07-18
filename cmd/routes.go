package main

import (
	"golang-forum-backend/controllers"
	"net/http"
)

func (a *application) routes() http.Handler {
	mux := http.NewServeMux()
	//mux.Use(middleware.RequestID)
	//mux.Use(middleware.RealIP)
	//mux.Use(middleware.Recoverer)
	//
	//if a.debug {
	//	mux.Use(middleware.Logger)
	//}
	mux.HandleFunc("/auth/signup", controllers.SignUp)
	mux.HandleFunc("/auth/google/signup", controllers.SignUpWithGoogle)
	mux.HandleFunc("/auth/google/callback", controllers.SignUpWithGoogleCallback)
	mux.HandleFunc("/auth/login", controllers.SignIn)
	mux.HandleFunc("/auth/github/login", controllers.HandleGitHubLogin)
	mux.HandleFunc("/auth/github/callback", controllers.HandleGitHubCallback)

	return mux
}
