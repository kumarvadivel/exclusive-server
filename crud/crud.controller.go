package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"go-server/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/addTodo", addTodo).Methods("POST")
	r.HandleFunc("/getTodos", getTodos).Methods("GET")
	r.HandleFunc("/getTodo/{id}", getTodo).Methods("GET")
}

func addTodo(res http.ResponseWriter, req *http.Request) {
	//log.Fatal("route Name:addTodo")
	res.Write([]byte(`helhshphf`))
}

func getTodos(res http.ResponseWriter, req *http.Request) {
	client := db.Client
	fmt.Print(client)
	var games []bson.M
	cursor, err := client.Database("nykaa").Collection("games").Find(context.TODO(), bson.M{})
	if err != nil {

	}

	if err = cursor.All(context.TODO(), &games); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(res).Encode(games)
}

func getTodo(res http.ResponseWriter, req *http.Request) {

}
