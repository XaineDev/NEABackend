package api

import (
	"NEABackend/src/database"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterFunction(writer http.ResponseWriter, request *http.Request) {

	registerrequest := RegisterRequest{
		Username: "",
		Password: "",
	}

	// read the request body into a buffer
	bodyBuffer := make([]byte, request.ContentLength)
	read, err := request.Body.Read(bodyBuffer)
	if err != nil && err != io.EOF {
		// return 500 if there is an error reading the body
		log.Println("Error reading request body: " + err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// parse the json from the body buffer into the struct
	err = json.Unmarshal(bodyBuffer[:read], &registerrequest)
	if err != nil {
		log.Println("Error parsing request body: " + err.Error())
		// return 500 if there is an error parsing the body
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(`{"success": false, "error": "invalid json"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// check to make sure that the username or password is not empty
	if registerrequest.Username == "" || registerrequest.Password == "" {
		// return 400 if the username or password is empty
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write([]byte(`{"success": false, "error": "invalid username or password"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// check if username is taken
	// if username is taken, return 400

	_, err = database.DatabaseConnection.GetUserByUsername(registerrequest.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("Error getting user: " + err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
			if err != nil {
				log.Println("Error writing response: " + err.Error())
			}
			return
		}
		// username is not taken, continue
	} else {
		// username is taken
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write([]byte(`{"success": false, "error": "username is taken"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// create user
	err = database.DatabaseConnection.CreateUser(registerrequest.Username, registerrequest.Password)
	if err != nil {
		log.Println("Error creating user: " + err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// get user
	user, err := database.DatabaseConnection.GetUserByUsername(registerrequest.Username)
	if err != nil {
		return
	}
	err = database.LogAction(user.ID, "register", "")
	if err != nil {
		return
	}

	// return 200
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(`{"success": true}`))
	if err != nil {
		log.Println("Error writing response: " + err.Error())
	}

}
