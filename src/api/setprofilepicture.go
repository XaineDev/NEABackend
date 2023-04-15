package api

import (
	"NEABackend/src/database"
	"NEABackend/src/util"
	"io"
	"log"
	"net/http"
)

type SetProfilePictureRequest struct {
	ProfileUsername string `json:"profile_username"` // username of the user to set the profile picture of
	ProfileID       string `json:"profile_id"`       // users id to verify request authenticity
	ProfileToken    string `json:"profile_token"`    // users token to verify request authenticity

	ProfilePicture string `json:"profile_picture"` // base 64 encoded image
}

type SetProfilePictureResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func SetProfilePictureFunction(writer http.ResponseWriter, request *http.Request) {

	responseStruct := SetProfilePictureResponse{
		Success: false,
		Error:   "unknown error",
	}

	requestStruct := SetProfilePictureRequest{
		ProfileUsername: "",
		ProfileID:       "",
		ProfileToken:    "",
		ProfilePicture:  "",
	}

	// read the request body into a buffer
	bodyBuffer := make([]byte, request.ContentLength)
	read, err := request.Body.Read(bodyBuffer)
	if err != nil && err != io.EOF {
		log.Println("Error reading request body: " + err.Error())
		responseStruct.Error = "internal server error"
		writer.WriteHeader(http.StatusInternalServerError)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}

		return
	}

	// parse the json from the body buffer into the struct
	err = util.ParseJson(bodyBuffer[:read], &requestStruct)
	if err != nil {
		log.Println("Error parsing request body: " + err.Error())
		responseStruct.Error = "invalid json"
		writer.WriteHeader(http.StatusInternalServerError)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}

		return
	}

	// check to make sure that the fields are not empty
	if requestStruct.ProfileUsername == "" || requestStruct.ProfileID == "" || requestStruct.ProfileToken == "" || requestStruct.ProfilePicture == "" {
		responseStruct.Error = "invalid json"
		writer.WriteHeader(http.StatusBadRequest)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}

		return
	}

	// check to make sure that the user is authenticated
	if user := database.DatabaseConnection.ValidateRequest(requestStruct.ProfileToken, requestStruct.ProfileID, requestStruct.ProfileUsername); user == nil {
		responseStruct.Error = "unauthorized"
		writer.WriteHeader(http.StatusUnauthorized)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}

		return
	} else {
		err = util.SaveImageFromBase64String(requestStruct.ProfilePicture, user.ID)
		if err != nil {
			log.Println("Error saving image: " + err.Error())
			responseStruct.Error = "internal server error"
			writer.WriteHeader(http.StatusInternalServerError)
			err = util.RespondWithJson(writer, responseStruct)
			if err != nil {
				log.Println("Error writing response: " + err.Error())
			}

			return
		}
	}

	// respond with success
	responseStruct.Success = true
	responseStruct.Error = ""

	err = util.RespondWithJson(writer, responseStruct)
	if err != nil {
		log.Println("Error writing response: " + err.Error())
	}

}
