package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Album struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	AlbumTitle  string             `json:"albumTitle"`
	Artist      string             `json:"artist"`
	ReleaseYear string             `json:"releaseYear"`
	Tracks      []Track            `json:"tracks"`
}

type Track struct {
	TrackTitle      string `json:"trackTitle"`
	DurationSeconds int    `json:"durationSeconds"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}

var client *mongo.Client

func ConnectToDB() {
	// connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	fmt.Println("Connected to MongoDB!")
}

// ErrorResponse : error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// GetError : This is a function to prepare the error model.
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

func createUser(w http.ResponseWriter, r *http.Request) {
	// get email and password from request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Bad request")
		w.WriteHeader(500)
	}

	// connect to users collection
	collection := client.Database("restAPI").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// insert user
	collection.InsertOne(ctx, bson.D{
		{"email", user.Email},
		{"password", user.Password},
	})

	if err != nil {
		// user not created
		log.Printf("Bad request")
		w.WriteHeader(500)
		return
	}

	fmt.Fprintf(w, "User created! %+v", user.Email)

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// create users array
	var users []User

	collection := client.Database("restAPI").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// pass empty filter to get all data
	cursor, err := collection.Find(ctx, bson.D{})

	if err != nil {
		GetError(err, w)
		return
	}

	// close cursor once finished
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user User

		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		// add album to our array
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(users)
}

func loginUser(w http.ResponseWriter, r *http.Request) {

	// get email and password from request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Unable to authenticate")
		w.WriteHeader(401)
		return
	}

	// get user from mongoDB by email address to see if user exists
	w.Header().Set("Content-Type", "application/json")

	// create dbUser variable for retrieving using from DB
	var dbUser User

	// connect to users collection
	collection := client.Database("restAPI").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get user from DB
	filter := bson.D{{"email", user.Email}}
	if err = collection.FindOne(ctx, filter).Decode(&dbUser); err != nil {
		// user not found or issue finding user
		log.Printf("Unable to authenticate")
		w.WriteHeader(401)
		return
	}

	if dbUser.Email == "" {
		// user not found
		log.Printf("Unable to authenticate")
		w.WriteHeader(401)
		return
	}

	isMatch := (user.Email == dbUser.Email && user.Password == dbUser.Password)

	if !isMatch {
		// user not matched
		log.Printf("Unable to authenticate")
		w.WriteHeader(401)
		return
	}

	fmt.Fprintf(w, "User is a match! %+v", dbUser)

}

func createAlbums(w http.ResponseWriter, r *http.Request) {
	// get array of albums from request body
	var albums []interface{}
	err := json.NewDecoder(r.Body).Decode(&albums)
	if err != nil {
		log.Printf("Unable to parse album data")
		w.WriteHeader(500)
		return
	}

	// connect to albums collection
	collection := client.Database("restAPI").Collection("albums")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// insert albums
	_, err = collection.InsertMany(ctx, albums)
	if err != nil {
		log.Printf("Unable to insert album data")
		w.WriteHeader(500)
		return
	}

	fmt.Fprintf(w, "Albums created! %+v", albums)
}

func updateAlbum(w http.ResponseWriter, r *http.Request) {
	// get updated album info from request body
	var album interface{}
	err := json.NewDecoder(r.Body).Decode(&album)
	if err != nil {
		log.Printf("Unable to parse album data")
		w.WriteHeader(500)
		return
	}

	// get album id from header
	vars := mux.Vars(r)
	key := vars["id"]

	id, _ := primitive.ObjectIDFromHex(key)

	// connect to albums collection
	collection := client.Database("restAPI").Collection("albums")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// update album
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{{"$set", album}},
	)

	if err != nil {
		log.Printf("Unable to update album data")
		w.WriteHeader(500)
		return
	}

	fmt.Fprintf(w, "Album updated!")
}

func deleteAlbum(w http.ResponseWriter, r *http.Request) {
	// get album id from header
	vars := mux.Vars(r)
	key := vars["id"]

	id, _ := primitive.ObjectIDFromHex(key)

	// connect to albums collection
	collection := client.Database("restAPI").Collection("albums")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// update album
	_, err := collection.DeleteOne(
		ctx,
		bson.M{"_id": id},
	)

	if err != nil {
		log.Printf("Unable to delete album data")
		w.WriteHeader(500)
		return
	}

	fmt.Fprintf(w, "Album deleted!")
}

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

// middleware handler for http handlers, enabling JWT authentication
func middleware(next http.Handler) http.Handler {
	SECRETKEY := "DAN_GOLANG_REST_API"
	// return our http handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the token from the authorisation header
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		// err if malformed
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		} else {
			// all ok so extract the token and try to parse it - signing method should be HMAC
			jwToken := authHeader[1]
			token, err := jwt.Parse(jwToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					// err if signing method is unexpected
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(SECRETKEY), nil
			})

			// get claims from the token - these are values encoded into the JWT
			// if decoding fails, the token has become corrupt or been tampered with
			// so return unauthorised
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), "props", claims)
				// Access context values in handlers - e.g.
				// props, _ := r.Context().Value("props").(jwt.MapClaims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
