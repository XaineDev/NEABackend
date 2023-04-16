package api

import (
	"NEABackend/src/database"
	"NEABackend/src/util"
	"io"
	"log"
	"net/http"
)

type UnClaimBookRequest struct {
	BookId    int          `json:"book_id"`
	BookTitle string       `json:"book_title"`
	User      *RequestUser `json:"user"`
}

type UnClaimBookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func UnclaimBookFunction(writer http.ResponseWriter, request *http.Request) {

	response := UnClaimBookResponse{
		Success: false,
		Message: "",
		Error:   "unknown error",
	}

	bodyBuffer := make([]byte, request.ContentLength)
	_, err := request.Body.Read(bodyBuffer)
	if err != nil && err != io.EOF {
		// return 500 if there is an error reading the body
		log.Println("Error reading body: " + err.Error())
		response.Error = "internal server error"
		writer.WriteHeader(http.StatusInternalServerError)
		err = util.RespondWithJson(writer, response)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	requestStruct := UnClaimBookRequest{
		BookId:    0,
		BookTitle: "",
		User:      nil,
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

	// check user is valid
	user, err := database.DatabaseConnection.GetUserByUsername(requestStruct.User.Username)
	if err != nil || user == nil {
		response.Error = "invalid user"
		writer.WriteHeader(http.StatusBadRequest)
		err = util.RespondWithJson(writer, response)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// check user token is valid
	if user.Token != requestStruct.User.Token {
		response.Error = "invalid token"
		writer.WriteHeader(http.StatusBadRequest)
		err = util.RespondWithJson(writer, response)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// check book exists
	book, err := database.DatabaseConnection.GetBookFromId(requestStruct.BookId)

	if err != nil || book == nil {
		response.Error = "invalid book"
		writer.WriteHeader(http.StatusBadRequest)
		err = util.RespondWithJson(writer, response)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// check book is claimed
	if book.CurrentOwner != user.ID {
		response.Error = "you do not own this book"
		writer.WriteHeader(http.StatusBadRequest)
		err = util.RespondWithJson(writer, response)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// unclaim book
	err = database.DatabaseConnection.UnclaimBook(book)
	if err != nil {
		log.Println("Error unclaiming book: " + err.Error())
		response.Error = "internal server error"
		writer.WriteHeader(http.StatusInternalServerError)
		err = util.RespondWithJson(writer, response)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// return success
	response.Success = true
	response.Message = "book unclaimed"
	response.Error = ""
	writer.WriteHeader(http.StatusOK)
	err = util.RespondWithJson(writer, response)

}
