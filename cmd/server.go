package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"net/http"
	"time"
)

func (a *application) listenAndServe() error {
	host := fmt.Sprintf("%s:%s", a.server.host, a.server.port)

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:8081"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
	)

	srv := http.Server{
		Handler:     corsHandler(a.routes()),
		Addr:        host,
		ReadTimeout: 300 * time.Second,
	}

	a.infoLog.Printf("Server listening on :%s", host)
	return srv.ListenAndServe()
}
