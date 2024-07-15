package main

import (
	"github.com/alexedwards/scs/v2"
	"golang-forum-backend/utils"
	"log"
	"os"
)

type application struct {
	appName string
	server  server
	debug   bool
	errLog  *log.Logger
	infoLog *log.Logger
	session *scs.SessionManager
}

type server struct {
	host string
	port string
	url  string
}

func main() {

	utils.InitDB()

	server := server{
		host: "localhost",
		port: "8082",
		url:  "http://localhost:8082",
	}

	app := &application{
		appName: "Golang Forum API",
		server:  server,
		debug:   true,
		errLog:  log.New(os.Stderr, "ERROR\t", log.Ltime|log.Ldate|log.Llongfile),
		infoLog: log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate|log.Lshortfile),
	}

	if err := app.listenAndServe(); err != nil {
		log.Fatal(err)
	}
}
