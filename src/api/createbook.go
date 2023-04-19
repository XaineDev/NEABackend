package api

import (
	"NEABackend/src/database"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type CreateBookRequest struct {
	AuthToken string `json:"authentication"`
	Title     string `json:"title"`
	Author    string `json:"author"`
}

type CreateBookResponse struct {
	Success bool          `json:"success"`
	Book    database.Book `json:"book,omitempty"`
}

func CreateBookFunction(writer http.ResponseWriter, request *http.Request) {
	// get the struct from the request body
	// check if the token is an admin token
	var createBookRequest CreateBookRequest
	// parse the json from the body buffer into the struct
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
	err = json.Unmarshal(bodyBuffer[:read], &createBookRequest)

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

	// check to make sure that the title and author is not empty
	if createBookRequest.Title == "" || createBookRequest.Author == "" {
		// return 400 if the title or author is empty
		log.Println("Title or Author not set")
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write([]byte(`{"success": false, "error": "invalid title or author"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	book := database.Book{
		Title:  createBookRequest.Title,
		Author: createBookRequest.Author,
	}

	// check if the token is an admin token
	validToken, err := database.DatabaseConnection.IsAdminToken(createBookRequest.AuthToken)

	if err != nil || !validToken {
		if err != nil {
			log.Println("Error checking token: " + err.Error())
		}
		// return 401 if the token is not an admin token
		log.Println("Invalid token")
		writer.WriteHeader(http.StatusUnauthorized)
		_, err = writer.Write([]byte(`{"success": false, "error": "invalid token"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	// add the book to the database
	err = database.DatabaseConnection.CreateBook(&book)
	if err != nil {
		// return 500 if there is an error creating the book
		log.Println("Error creating book: " + err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	response := CreateBookResponse{
		Success: true,
		Book:    book,
	}

	// return 200 and the struct if the book was created successfully
	responseBytes, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshaling response: " + err.Error())
		_, err = writer.Write([]byte(`{"success": false, "error": "internal server error"}`))
		if err != nil {
			log.Println("Error writing response: " + err.Error())
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(responseBytes)
	if err != nil {
		log.Println("Error writing response: " + err.Error())
	}

	user, err := database.DatabaseConnection.GetUserByToken(createBookRequest.AuthToken)
	if err != nil {
		log.Println("Error getting user: " + err.Error())
	}
	err = database.LogAction(user.ID, "created book", "book: "+(createBookRequest.Title)+" by "+(createBookRequest.Author))
	if err != nil {
		log.Println("Error logging action: " + err.Error())
	}
}
