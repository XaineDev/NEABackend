package api

import (
	"NEABackend/src/database"
	"NEABackend/src/util"
	"database/sql"
	"log"
	"net/http"
)

type FetchLogsRequest struct {
	User       AdminAuthentication `json:"authentication"`        // user used for authenticating the request
	Amount     int                 `json:"amount"`                // amount of logs to fetch
	TargetUser int                 `json:"target_user,omitempty"` // target user to fetch logs for
	Page       int                 `json:"page"`                  // page of logs to start from
}

type FetchLogsResponse struct {
	Success bool                   `json:"success"`         // whether the request was successful
	Actions []*database.UserAction `json:"actions"`         // list of actions
	Error   string                 `json:"error,omitempty"` // error message if the request was not successful
}

func FetchLogsFunction(writer http.ResponseWriter, request *http.Request) {

	fetchLogsRequest := FetchLogsRequest{}
	err := util.ParseRequest(request, &fetchLogsRequest)
	if err != nil {
		// return 400 if there is an error parsing the request
		writer.WriteHeader(http.StatusBadRequest)
		err := util.RespondWithJson(writer, FetchLogsResponse{
			Success: false,
			Actions: []*database.UserAction{},
			Error:   err.Error(),
		})
		if err != nil {
			log.Println("Error writing response: " + err.Error())
			return
		}
		return
	}
	if fetchLogsRequest.Amount == 0 {
		fetchLogsRequest.Amount = 25
	}

	user, err := database.DatabaseConnection.GetUserByToken(fetchLogsRequest.User.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			// return 401 if the user is not found
			writer.WriteHeader(http.StatusUnauthorized)
			err := util.RespondWithJson(writer, FetchLogsResponse{
				Success: false,
				Actions: []*database.UserAction{},
				Error:   "invalid token",
			})
			if err != nil {
				log.Println("Error writing response: " + err.Error())
				return
			}
			return
		}
	}

	if user == nil {
		// return 401 if the user is not found
		writer.WriteHeader(http.StatusUnauthorized)
		err := util.RespondWithJson(writer, FetchLogsResponse{
			Success: false,
			Actions: []*database.UserAction{},
			Error:   "invalid token",
		})
		if err != nil {
			log.Println("Error writing response: " + err.Error())
			return
		}
		return
	}

	if user.Username != fetchLogsRequest.User.Username || user.IsAdmin == false {
		// return 401 if the user is not found
		writer.WriteHeader(http.StatusUnauthorized)
		err := util.RespondWithJson(writer, FetchLogsResponse{
			Success: false,
			Actions: []*database.UserAction{},
			Error:   "invalid token",
		})
		if err != nil {
			log.Println("Error writing response: " + err.Error())
			return
		}
		return
	}

	if fetchLogsRequest.TargetUser == 0 {
		fetchAll(writer, fetchLogsRequest)
	} else {
		fetchTarget(writer, fetchLogsRequest)
	}

}

func fetchAll(writer http.ResponseWriter, fetchLogsRequest FetchLogsRequest) {
	actions, err := database.FetchActions(fetchLogsRequest.Amount, fetchLogsRequest.Page*fetchLogsRequest.Amount)
	if err != nil {
		// return 500 if there is an error getting the actions
		writer.WriteHeader(http.StatusInternalServerError)
		responseErr := util.RespondWithJson(writer, FetchLogsResponse{
			Success: false,
			Actions: []*database.UserAction{},
			Error:   "internal server error",
		})
		if responseErr != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
	err = util.RespondWithJson(writer, FetchLogsResponse{
		Success: true,
		Actions: actions,
		Error:   "",
	})
	if err != nil {
		log.Println("Error writing response: " + err.Error())
		return
	}
}

func fetchTarget(writer http.ResponseWriter, fetchLogsRequest FetchLogsRequest) {
	actions, err := database.FetchActionsFromUser(fetchLogsRequest.TargetUser, fetchLogsRequest.Amount, fetchLogsRequest.Page*fetchLogsRequest.Amount)
	if err != nil {
		// return 500 if there is an error getting the actions
		writer.WriteHeader(http.StatusInternalServerError)
		responseErr := util.RespondWithJson(writer, FetchLogsResponse{
			Success: false,
			Actions: []*database.UserAction{},
			Error:   "internal server error",
		})
		if responseErr != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
	err = util.RespondWithJson(writer, FetchLogsResponse{
		Success: true,
		Actions: actions,
		Error:   "",
	})
	if err != nil {
		log.Println("Error writing response: " + err.Error())
		return
	}
}
