package main

import (
	"BDMySQL/docs"
	"BDMySQL/storage"
	"database/sql"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	// запускаем логирование в файл app.log
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)

	// swagger
	docs.SwaggerInfo.Title = "Swagger Example API"

	// открывем БД
	db, err := sql.Open("mysql", "root:password@/userdb")

	s := &storage.Server{Database: db}

	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", s.IndexHandler).Methods("GET")
	router.HandleFunc("/create", s.CreateHandler).Methods("POST")
	router.HandleFunc("/user/{id}", s.EditPage).Methods("GET")
	router.HandleFunc("/edit/{id}", s.EditHandler).Methods("PATCH")
	router.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)
	http.Handle("/", router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
