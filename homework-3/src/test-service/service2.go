package main

import (
	"net/http"
	"github.com/gorilla/mux"
)


func main() {
    
	router := mux.NewRouter()

	router.HandleFunc("/posts", getAll).Methods("GET")

	router.HandleFunc("/posts/{id}", get).Methods("GET")

	http.ListenAndServe(":8081", router)
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Dobar dan!"))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Dobar dan!"))
}

