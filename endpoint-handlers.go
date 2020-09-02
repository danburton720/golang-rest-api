package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAlbums(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// create albums array
	var albums []Album

	collection := client.Database("restAPI").Collection("albums")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// pass empty filter to get all data
	cursor, err := collection.Find(ctx, bson.D{})

	if err != nil {
		w.WriteHeader(500)
		log.Printf("Bad request")
		w.WriteHeader(500)
		return
	}

	// close cursor once finished
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var album Album

		err := cursor.Decode(&album)
		if err != nil {
			log.Printf("Bad request")
			w.WriteHeader(500)
			return
		}

		// add album to our array
		albums = append(albums, album)
	}

	if err := cursor.Err(); err != nil {
		w.WriteHeader(500)
		log.Printf("Bad request")
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(albums)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// get email and password from request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Bad request")
		w.WriteHeader(500)
		return
	}

	// connect to users collection
	collection := client.Database("restAPI").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// insert user
	_, err = collection.InsertOne(ctx, bson.D{
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

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// create users array
	var users []User

	collection := client.Database("restAPI").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// pass empty filter to get all data
	cursor, err := collection.Find(ctx, bson.D{})

	if err != nil {
		log.Printf("Bad request")
		w.WriteHeader(500)
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
		log.Printf("Bad request")
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

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

	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	claims := TokenClaims{
		user.ID,
		user.Email,
		jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	tokenString, error := token.SignedString([]byte(goDotEnvVariable("SECRET_KEY")))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = user.Email

	json.NewEncoder(w).Encode(resp)
}

func CreateAlbums(w http.ResponseWriter, r *http.Request) {
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

func UpdateAlbum(w http.ResponseWriter, r *http.Request) {
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
	errNoDocuments := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.D{{"$set", album}},
	)

	if errNoDocuments != nil {
		log.Printf("Unable to update album data")
		w.WriteHeader(500)
		return
	}

	fmt.Fprintf(w, "Album updated!")
}

func DeleteAlbum(w http.ResponseWriter, r *http.Request) {
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
