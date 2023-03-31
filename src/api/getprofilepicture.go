package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GetProfilePictureRequest struct {
	ProfileUsername string `json:"profile_username"`
}

type GetProfilePictureResponse struct {
}

var GetProfilePictureFunction = func(writer http.ResponseWriter, request *http.Request) {
	// parse request json into struct
	receivedRequest := GetProfilePictureRequest{
		ProfileUsername: "",
	}

	// read the request body into a buffer
	bodyBuffer := make([]byte, request.ContentLength)
	read, err := request.Body.Read(bodyBuffer)
	if err != nil && err != io.EOF {
		// return 500 if there is an error reading the body
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// parse the json from the body buffer into the struct
	err = json.Unmarshal(bodyBuffer[:read], &receivedRequest)
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

	// check to make sure that the username is not empty
	if receivedRequest.ProfileUsername == "" {
		// return 400 if the username is empty
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write([]byte(`{"success": false, "error": "invalid username"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// check to see if the user exists in the database

	// get the users profile picture from the data folder

}
