package api

import (
	"NEABackend/src/database"
	"NEABackend/src/util"
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool          `json:"success"`
	User    database.User `json:"user,omitempty"`
}

func LoginFunction(writer http.ResponseWriter, request *http.Request) {

	loginRequest := LoginRequest{
		Username: "",
		Password: "",
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
	err = json.Unmarshal(bodyBuffer[:read], &loginRequest)
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
	if loginRequest.Username == "" || loginRequest.Password == "" {
		// return 400 if the username or password is empty
		log.Println("Username or Password not set")
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write([]byte(`{"success": false, "error": "invalid username or password"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// validate the username and password in database
	// if valid, return 200 with a token
	// if invalid, return 401 with an error message

	user, err := database.DatabaseConnection.GetUserByUsername(loginRequest.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			// no user found
			writer.WriteHeader(http.StatusUnauthorized)
			_, err = writer.Write([]byte(`{"success": false, "error": "invalid username or password"}`))
			if err != nil {
				log.Println("Error writing response: " + err.Error())
			}
			return
		}
		log.Println("Error getting user from database: " + err.Error())
		// return 500 if there is an error reading the body
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// hash password given with the salt from the database
	// compare the hash with the hash in the database
	// if they match, return 200 with a token
	// if they don't match, return 401 with an error message

	checkPassword, err := util.HashPasswordWithSalt(loginRequest.Password, user.PasswordSalt)
	if err != nil {
		// return 500 if there is an error reading the body
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	if subtle.ConstantTimeCompare([]byte(checkPassword), []byte(user.Password)) == 1 {
		// passwords match

		// generate a token
		token, err := util.CreateToken()
		if err != nil {
			// return 500 if there is an error reading the body
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
			if err != nil {
				log.Println("Error writing response: " + err.Error())
			}
			return
		}

		// update the user's token in the database
		err = database.DatabaseConnection.UpdateUserToken(user, token)
		if err != nil {
			// return 500 if there is an error reading the body
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
			if err != nil {
				log.Println("Error writing response: " + err.Error())
			}
			return
		}

		// update user and hide password and salt from response
		user.Token = token
		user.Password = ""
		user.PasswordSalt = ""

		response := LoginResponse{
			Success: true,
			User:    *user,
		}

		// marshal the response into json
		responseJson, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshalling response: " + err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
			if err != nil {
				log.Println("Error writing response: " + err.Error())
			}
			return
		}

		err = database.LogAction(user.ID, "login", "")
		if err != nil {
			return
		}

		// return 200 with the token
		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(responseJson)
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	} else {
		// passwords don't match
		// return 401 with an error message
		writer.WriteHeader(http.StatusUnauthorized)
		_, err = writer.Write([]byte(`{"success": false, "error": "invalid username or password"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

}
