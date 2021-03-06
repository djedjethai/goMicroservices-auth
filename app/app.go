package app

import (
	"fmt"
	"github.com/djedjethai/bankingAuth/domain"
	"github.com/djedjethai/bankingAuth/service"
	"github.com/djedjethai/bankingLib/logger"
	_ "github.com/go-sql-driver/mysql"
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
	dom := domain.NewAuthRepository(dbClient)
	// add domain into service
	s := service.NewService(dom, domain.GetRolePermissions())
	// create the authHadler
	ah := authHandler{s}

	r := mux.NewRouter()
	r.HandleFunc("/auth/login", ah.login).Methods(http.MethodPost)
	r.HandleFunc("/auth/register", ah.addCustomer).Methods(http.MethodPost)
	r.HandleFunc("/auth/refresh", ah.refresh).Methods(http.MethodPost)
	r.HandleFunc("/auth/verify", ah.verify).Methods(http.MethodGet)
	address := os.Getenv("SERVER_ADDR")
	port := os.Getenv("SERVER_PORT")
	logger.Info(fmt.Sprintf("starting Oauth server on %s:%s", address, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), r))
}

func getDbClient() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbAddr := os.Getenv("DB_ADDR")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasswd, dbAddr, dbPort, dbName)
	client, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client
}
