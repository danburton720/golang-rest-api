package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectToDB() {
	// connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	fmt.Println("Connected to MongoDB!")
}

func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var header = r.Header.Get("Authorization")
		header = strings.TrimSpace(header)
		extractedToken := strings.Split(header, "Bearer ")
		extractedTokenFinal := extractedToken[1]

		if header == "" {
			// token is missing so return 403 Unauthorized
			log.Printf("Missing auth token 1")
			w.WriteHeader(403)
			return
		}

		// TODO - define token
		type TokenClaims struct {
			ID    primitive.ObjectID `bson:"_id,omitempty"`
			Email string             `json:"Email"`
			jwt.StandardClaims
		}

		token, _ := jwt.ParseWithClaims(extractedTokenFinal, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("DAN_GOLANG_REST_API"), nil
		})

		tokenValid := false
		if _, ok := token.Claims.(*TokenClaims); ok && token.Valid {
			tokenValid = true
		}

		if tokenValid {
			ctx := context.WithValue(r.Context(), "user", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			log.Printf("Missing auth token 2")
			w.WriteHeader(403)
			return
		}
	})

}

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
