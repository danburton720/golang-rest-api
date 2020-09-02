# golang-rest-api

## requires
* go
* go.mongodb.org/mongo-driver 
* github.com/gorilla/mux
* github.com/dgrijalva/jwt-go
* github.com/joho/godotenv
* golang.org/x/crypto/bcrypt

mongoDB runs in local host port 27017

run powershell script, e.g. /Users/user/mongodb/bin/mongod.exe --dbpath=/Users/user/mongodb-data

Go API runs in localhost on port 8081 - use 'go *.go' to start

Bearer token required to hit /auth endpoints e.g. /users and /albums and /albums/{id}

This can be retrieved when logging in i.e. /users/login endpoint passing user name and password, then use token in requests to those endpoints

## current restrictions / TODO:
* object ids from mongo should be returned with collection representations