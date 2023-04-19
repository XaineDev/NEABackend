package api

import (
	"NEABackend/src/database"
	"NEABackend/src/util"
	"net/http"
	"strconv"
)

type AdminAuthentication struct {
	ID       int    `json:"id"`
	Token    string `json:"token"`
	Username string `json:"username"`
}

type UnClaimBookRequest struct {
	BookId         int                 `json:"book_id"`
	BookTitle      string              `json:"book_title"`
	CurrentOwner   int                 `json:"current_owner"`
	Authentication AdminAuthentication `json:"authentication"`
}

type UnClaimBookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func ForceUnclaimFunction(writer http.ResponseWriter, request *http.Request) {

	var requestStruct UnClaimBookRequest

	err := util.ParseRequest(request, &requestStruct)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_ = util.RespondWithJson(writer, UnClaimBookResponse{
			Success: false,
			Message: "",
			Error:   "invalid json",
		})
		return
	}

	// validate the authentication
	user, err := database.DatabaseConnection.GetUserByUsername(requestStruct.Authentication.Username)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_ = util.RespondWithJson(writer, UnClaimBookResponse{
			Success: false,
			Message: "",
			Error:   "internal server error",
		})
		return
	}

	if user.ID != requestStruct.Authentication.ID || user.Token != requestStruct.Authentication.Token {
		writer.WriteHeader(http.StatusUnauthorized)
		_ = util.RespondWithJson(writer, UnClaimBookResponse{
			Success: false,
			Message: "",
			Error:   "unauthorized",
		})
		return
	}

	// validate the book
	book, err := database.DatabaseConnection.GetBookFromId(requestStruct.BookId)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_ = util.RespondWithJson(writer, UnClaimBookResponse{
			Success: false,
			Message: "",
			Error:   "internal server error",
		})
		return
	}

	if book == nil {
		writer.WriteHeader(http.StatusBadRequest)
		_ = util.RespondWithJson(writer, UnClaimBookResponse{
			Success: false,
			Message: "",
			Error:   "book does not exist",
		})
		return
	}

	if book.Title != requestStruct.BookTitle {
		writer.WriteHeader(http.StatusBadRequest)
		_ = util.RespondWithJson(writer, UnClaimBookResponse{
			Success: false,
			Message: "",
			Error:   "book title does not match",
		})
		return
	}

	if book.CurrentOwner != requestStruct.CurrentOwner {
		writer.WriteHeader(http.StatusBadRequest)
		_ = util.RespondWithJson(writer, UnClaimBookResponse{
			Success: false,
			Message: "",
			Error:   "book is not currently claimed by the user",
		})
		return
	}

	// unclaim the book
	err = database.DatabaseConnection.UnclaimBook(book)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_ = util.RespondWithJson(writer, UnClaimBookResponse{
			Success: false,
			Message: "",
			Error:   "internal server error",
		})
		return
	}

	_ = database.LogAction(user.ID, "force unclaim",
		"book: "+strconv.Itoa(book.ID)+" owner: "+strconv.Itoa(requestStruct.CurrentOwner)+"")

	// return success
	writer.WriteHeader(http.StatusOK)
	_ = util.RespondWithJson(writer, UnClaimBookResponse{
		Success: true,
		Message: "book unclaimed",
		Error:   "",
	})
}
