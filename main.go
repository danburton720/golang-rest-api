package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

func getAlbums(w http.ResponseWriter, r *http.Request) {

	Albums := []Album{
		Album{
			AlbumTitle:  "City of Evil",
			Artist:      "Avenged Sevenfold",
			ReleaseYear: "2005",
			Tracks: []Track{
				{
					TrackTitle:      "Beast and the Harlot",
					DurationSeconds: 344,
				},
				{
					TrackTitle:      "Burn it Down",
					DurationSeconds: 299,
				},
			},
		},
		{
			AlbumTitle:  "Emperor of Sand",
			Artist:      "Mastodon",
			ReleaseYear: "2017",
			Tracks: []Track{
				{
					TrackTitle:      "Sultan's Curse",
					DurationSeconds: 250,
				},
				{
					TrackTitle:      "Show Yourself",
					DurationSeconds: 183,
				},
			},
		},
	}

	fmt.Println("Endpoint Hit: Get Albums Endpoint")
	json.NewEncoder(w).Encode(Albums)
}

func getTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Endpoint Hit: Get Test Endpoint")
}

func testPostAlbums(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test POST endpoint worked")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/albums", getAlbums).Methods("GET")
	myRouter.HandleFunc("/test", getTest).Methods("GET")
	myRouter.HandleFunc("/albums", testPostAlbums).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	handleRequests()
}
