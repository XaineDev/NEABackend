package api

import (
	"NEABackend/src/database"
	"NEABackend/src/util"
	"io"
	"log"
	"net/http"
)

type GetBooksRequest struct {
	User *RequestUser `json:"user"`
}

type GetBooksResponse struct {
	Success bool             `json:"success"`
	Books   []*database.Book `json:"books"`
	Error   string           `json:"error,omitempty"`
}

func GetBooksFunction(writer http.ResponseWriter, request *http.Request) {

	response := GetBooksResponse{
		Success: false,
		Books:   nil,
		Error:   "unknown error",
	}

	bodyBuffer := make([]byte, request.ContentLength)
	_, err := request.Body.Read(bodyBuffer)
	if err != nil && err != io.EOF {
		// return 500 if there is an error reading the body
		response.Error = "internal server error"
		writer.WriteHeader(http.StatusInternalServerError)
		err = util.RespondWithJson(writer, response)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	requestStruct := GetBooksRequest{
		User: nil,
	}

	err = util.ParseJson(bodyBuffer, &requestStruct)
	if err != nil {
		// return 400 if there is an error parsing the json
		response.Error = "invalid json"
		writer.WriteHeader(http.StatusBadRequest)
		err = util.RespondWithJson(writer, response)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	user, err := database.DatabaseConnection.GetUserByUsername(requestStruct.User.Username)
	if err != nil {
		// return 500 if there is an error getting the user
		log.Println("Error getting user: " + err.Error())
		response.Error = "internal server error"
		writer.WriteHeader(http.StatusInternalServerError)
		err = util.RespondWithJson(writer, response)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	if user == nil || user.Token != requestStruct.User.Token {
		// return 401 if the user is not found
		response.Error = "invalid user"
		writer.WriteHeader(http.StatusUnauthorized)
		err = util.RespondWithJson(writer, response)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	books, err := database.DatabaseConnection.GetBooks()

	if err != nil {
		// return 500 if there is an error getting the books
		log.Println("Error getting books: " + err.Error())
		response.Error = "internal server error"
		writer.WriteHeader(http.StatusInternalServerError)
		err = util.RespondWithJson(writer, response)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	response.Success = true
	response.Books = books
	response.Error = ""
	writer.WriteHeader(http.StatusOK)
	err = util.RespondWithJson(writer, response)

}
