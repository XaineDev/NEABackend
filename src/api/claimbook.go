package api

import (
	"NEABackend/src/database"
	"NEABackend/src/util"
	"io"
	"log"
	"net/http"
)

type ClaimBookRequest struct {
	BookId    int          `json:"book_id"`
	BookTitle string       `json:"book_title"`
	User      *RequestUser `json:"user"`
}

type RequestUser struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type ClaimBookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func ClaimBookFunction(writer http.ResponseWriter, request *http.Request) {

	responseStruct := ClaimBookResponse{
		Success: false,
		Message: "",
		Error:   "unknown error",
	}

	bodyBuffer := make([]byte, request.ContentLength)
	_, err := request.Body.Read(bodyBuffer)
	if err != nil && err != io.EOF {
		// return 500 if there is an error reading the body
		responseStruct.Error = "internal server error"
		writer.WriteHeader(http.StatusInternalServerError)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	requestStruct := ClaimBookRequest{
		BookId:    0,
		BookTitle: "",
		User:      nil,
	}

	err = util.ParseJson(bodyBuffer, &requestStruct)
	if err != nil {
		// return 400 if there is an error parsing the json
		responseStruct.Error = "invalid json"
		writer.WriteHeader(http.StatusBadRequest)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// check user is valid
	user, err := database.DatabaseConnection.GetUserByUsername(requestStruct.User.Username)
	if err != nil {
		log.Println("Error getting user from username: " + err.Error())
		responseStruct.Error = "internal server error"
		writer.WriteHeader(http.StatusInternalServerError)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	if user == nil {
		responseStruct.Error = "user not found"
		writer.WriteHeader(http.StatusBadRequest)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	if user.Token != requestStruct.User.Token {
		responseStruct.Error = "invalid token"
		writer.WriteHeader(http.StatusBadRequest)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// check if book is available
	book, err := database.DatabaseConnection.GetBookFromId(requestStruct.BookId)
	if err != nil {
		log.Println("Error getting book from id: " + err.Error())
		responseStruct.Error = "internal server error"
		writer.WriteHeader(http.StatusInternalServerError)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	if book == nil {
		responseStruct.Error = "book not found"
		writer.WriteHeader(http.StatusBadRequest)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	if !book.IsAvailable() {
		responseStruct.Error = "book not available"
		writer.WriteHeader(http.StatusBadRequest)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// update book with new owner
	book.CurrentOwner = user.Username
	err = book.Update()
	if err != nil {
		log.Println("Error updating book: " + err.Error())
		responseStruct.Error = "failed to update book"
		writer.WriteHeader(http.StatusInternalServerError)
		err = util.RespondWithJson(writer, responseStruct)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	responseStruct.Success = true
	responseStruct.Message = "book claimed successfully"
	responseStruct.Error = ""
	writer.WriteHeader(http.StatusOK)
	err = util.RespondWithJson(writer, responseStruct)
}
