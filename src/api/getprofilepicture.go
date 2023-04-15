package api

import (
	"NEABackend/src/database"
	"NEABackend/src/util"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type GetProfilePictureRequest struct {
	ProfileUsername string `json:"profile_username"`
}

type GetProfilePictureResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
	ImageData string `json:"image_data,omitempty"`
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
	user, err := database.DatabaseConnection.GetUserByUsername(receivedRequest.ProfileUsername)
	if err != nil {
		// return 500 if there is an error getting the user from the database
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	if user == nil {
		// return 400 if the user does not exist
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write([]byte(`{"success": false, "error": "user not found"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// check if user has a profile picture
	// profile pictures are stored as the userid.png in the data/profiles folder

	imageData, err := util.GetImageAsBase64String("data/profiles/" + strconv.Itoa(user.ID) + ".png")
	if err != nil {
		if err == os.ErrNotExist {
			// return 400 if the image does not exist
			writer.WriteHeader(http.StatusBadRequest)
			_, err = writer.Write([]byte(`{"success": false, "error": "profile picture not found"}`))
			if err != nil {
				log.Println("Error writing response: " + err.Error())
			}
			return
		}
		// return 500 if there is an unknown error getting the image
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	responseStruct := GetProfilePictureResponse{
		Success:   true,
		Error:     "",
		ImageData: imageData,
	}

	// marshal the response struct into json
	responseJson, err := json.Marshal(responseStruct)
	if err != nil {
		// return 500 if there is an error marshaling the response struct
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// write the response
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(responseJson)
	if err != nil {
		log.Println("Error writing response: " + err.Error())
	}

}
