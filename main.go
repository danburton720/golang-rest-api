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
	r.Use(CommonMiddleware)

	r.HandleFunc("/", HomePage)
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users/login", LoginUser).Methods("POST")

	// Auth route
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(JwtVerify)

	s.HandleFunc("/users", GetUsers).Methods("GET")
	s.HandleFunc("/albums", GetAlbums).Methods("GET")
	s.HandleFunc("/albums", CreateAlbums).Methods("POST")
	s.HandleFunc("/albums/{id}", UpdateAlbum).Methods("PUT")
	s.HandleFunc("/albums/{id}", DeleteAlbum).Methods("DELETE")

	log.Fatal(http.ListenAndServe(goDotEnvVariable("PORT"), r))
}
