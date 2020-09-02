# golang-rest-api

## requires
* go,
* go.mongodb.org/mongo-driver 
* github.com/gorilla/mux, 
* github.com/dgrijalva/jwt-go
* github.com/joho/godotenv

mongoDB runs in local host port 27017
run powershell script, e.g. /Users/<user>/mongodb/bin/mongod.exe --dbpath=/Users/<user>/mongodb-data

Go API runs in localhost on port 8081 - use 'go *.go' to start

Bearer token required to hit /auth endpoints e.g. /users and /albums and /albums/{id}
This can be retrieved when logging in i.e. /users/login endpoint passing user name and password

## current restrictions / TODO:
* passwords stored as plaintext (need to hash before storing)
* use more "bad request" responses rather than fatal errors
* custom errors to be implemented better in general