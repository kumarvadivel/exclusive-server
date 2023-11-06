package main

import (
	"fmt"
	"go-server/controllers"
	"go-server/crud"
	"go-server/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Startig Application...")
	db.MongoConnect()
	r := mux.NewRouter()
	crud.RegisterRoutes(r)
	controllers.RegisterRoutes(r)

	http.Handle("/", r)

	http.ListenAndServe(":3000", nil)
}
