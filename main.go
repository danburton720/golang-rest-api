package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	fmt.Println("Connecting to DB")
	ConnectToDB()

	r := mux.NewRouter().StrictSlash(true)
	// r.Use(middleware)
	r.HandleFunc("/", homePage)
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/login", loginUser).Methods("POST")
	r.HandleFunc("/albums", getAlbums).Methods("GET")
	r.HandleFunc("/albums", createAlbums).Methods("POST")
	r.HandleFunc("/albums/{id}", updateAlbum).Methods("PUT")
	r.HandleFunc("/albums/{id}", deleteAlbum).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8081", r))
}
