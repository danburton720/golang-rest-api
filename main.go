package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Album struct {
	AlbumTitle  string  `json:"albumTitle"`
	Artist      string  `json:"artist"`
	ReleaseYear string  `json:"releaseYear"`
	Tracks      []Track `json:"tracks"`
}

type Track struct {
	TrackTitle      string `json:"trackTitle"`
	DurationSeconds int    `json:"durationSeconds"`
}

var client *mongo.Client

func ConnectToDB() {

	// connect to MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	fmt.Println("Connected to MongoDB!")

}

// ErrorResponse : error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// GetError : This is a helper function to prepare the error model.
func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

func getAlbums(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// create albums array
	var albums []Album

	collection := client.Database("restAPI").Collection("albums")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// pass empty filter to get all data
	cursor, err := collection.Find(ctx, bson.D{})

	if err != nil {
		GetError(err, w)
		return
	}

	// close cursor once finished
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var album Album

		err := cursor.Decode(&album)
		if err != nil {
			log.Fatal(err)
		}

		// add album to our array
		albums = append(albums, album)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(albums)
}

func testPostAlbums(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test POST endpoint worked")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func main() {

	fmt.Println("Connecting to DB")
	ConnectToDB()

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", homePage)
	r.HandleFunc("/albums", getAlbums).Methods("GET")
	r.HandleFunc("/albums", testPostAlbums).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", r))
}
