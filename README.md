# golang-rest-api

## requires
* go,
* go.mongodb.org/mongo-driver 
* github.com/gorilla/mux, 
* github.com/dgrijalva/jwt-go

mongoDB runs in local host port 27017
run powershell script, e.g. /Users/<user>/mongodb/bin/mongod.exe --dbpath=/Users/<user>/mongodb-data

Go API runs in localhost on port 8081 - use 'go *.go' to start

## current restrictions / TODO:
* passwords stored as plaintext (need to hash before storing)
* middleware set up for JWT authentication but currently not utilised (WIP)
* user can "login" but as JWT authentication not complete, this doesn't actually do anything (you can still access the other endpoints)
* use more "bad request" responses rather than fatal errors
* custom errors to be implemented better in general
* refactor and split code into multiple files rather than all in one main.go