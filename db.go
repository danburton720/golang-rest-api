package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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
