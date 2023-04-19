package api

import (
	"NEABackend/src/database"
	"NEABackend/src/util"
	"log"
	"net/http"
)

type FetchUsersRequest struct {
	User   AdminAuthentication `json:"authentication"` // user used for authenticating the request
	Amount int                 `json:"amount"`         // amount of users to fetch
	Page   int                 `json:"page"`           // page of users to start from
}

type FetchUsersResponse struct {
	Success bool            `json:"success"`         // whether the request was successful
	Users   []database.User `json:"users"`           // list of actions
	Error   string          `json:"error,omitempty"` // error message if the request was not successful
}

func FetchUsersFunction(writer http.ResponseWriter, request *http.Request) {

	var fetchUsersRequest FetchUsersRequest
	err := util.ParseRequest(request, &fetchUsersRequest)
	if err != nil {
		// return 400 if there is an error parsing the request
		writer.WriteHeader(http.StatusBadRequest)
		err := util.RespondWithJson(writer, FetchUsersResponse{
			Success: false,
			Users:   []database.User{},
			Error:   err.Error(),
		})
		if err != nil {
			log.Println("Error writing response: " + err.Error())
			return
		}
		return
	}

	if fetchUsersRequest.Amount == 0 {
		fetchUsersRequest.Amount = 25
	}

	user, err := database.DatabaseConnection.GetUserByToken(fetchUsersRequest.User.Token)
	if err != nil {
		// return 401 if the user is not found
		writer.WriteHeader(http.StatusUnauthorized)
		err := util.RespondWithJson(writer, FetchUsersResponse{
			Success: false,
			Users:   []database.User{},
			Error:   "invalid token",
		})
		if err != nil {
			log.Println("Error writing response: " + err.Error())
			return
		}
		return
	}

	if !user.IsAdmin {
		// return 401 if the user is not an admin
		writer.WriteHeader(http.StatusUnauthorized)
		err := util.RespondWithJson(writer, FetchUsersResponse{
			Success: false,
			Users:   []database.User{},
			Error:   "not an admin",
		})
		if err != nil {
			log.Println("Error writing response: " + err.Error())
			return
		}
		return
	}

	users, err := database.DatabaseConnection.GetUsers(fetchUsersRequest.Amount, fetchUsersRequest.Page)
	if err != nil {
		// return 500 if there is an error fetching the users
		writer.WriteHeader(http.StatusInternalServerError)
		err := util.RespondWithJson(writer, FetchUsersResponse{
			Success: false,
			Users:   []database.User{},
			Error:   "error fetching users",
		})
		if err != nil {
			log.Println("Error writing response: " + err.Error())
			return
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
	err = util.RespondWithJson(writer, FetchUsersResponse{
		Success: true,
		Users:   users,
	})
	if err != nil {
		log.Println("Error writing response: " + err.Error())
		return
	}

}
