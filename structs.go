package main

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type TokenClaims struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	jwt.StandardClaims
}
