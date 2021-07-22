package app

import (
	"github.com/djedjethai/bankingAuth/logger"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"time"
)

func sanityCheck() {
	if os.Getenv("SERVER_ADDR") == "" ||
		os.Getenv("SERVER_PORT") == "" ||
		os.Getenv("DB_USER") == "" ||
		os.Getenv("DB_PASSWD") == "" ||
		os.Getenv("DB_ADDR") == "" ||
		os.Getenv("DB_PORT") == "" ||
		os.Getenv("DB_NAME") == "" {
		log.Fatal("Environment variables not define")
	}
}

func Start() {

	sanityCheck()

	dbClient := getDbClient()

	// add dbclient into domain

	// add domain into service

	// svc instance
	// ah := ...

	r := mux.NewRouter()
	r.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)
	// r.HandleFunc("/auth/register", ah.NotImplemented).Methods(http.MethodPost)
	// r.HandleFunc("/auth/refresh", ah.Refresh).Methods(http.MethodPost)
	// r.HandleFunc("/auth/verify", ah.Verify).Methods(http.MethodGet)

	address := os.Getenv("SERVER_ADDR")
	port := os.Getenv("SERVER_PORT")
	logger.Info(fmt.Sprintf("starting Oauth server on %s:%s", address, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("Listen on: %s:%s", address, port), NewRouter))
}

func getDbClient() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbAddr := os.Getenv("DB_ADDR")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dataSource := fmt.Printf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasswd, dbAddr, dbPort, dbName)
	client, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client
}
