package main

import (
	"NEABackend/src/api"
	"NEABackend/src/database"
	"NEABackend/src/util"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	BackendPort   = flag.String("port", "3000", "port to host the api on")
	router        chi.Router
	RateLimits    = map[string]time.Time{}
	RateLimitTime = flag.Int("ratelimit", 0, "ratelimit for api requests in milliseconds")
)

func init() {
	flag.Parse()
}

func init() {
	// initiate a router to easily handle requests
	router = chi.NewRouter()

	// create our custom middleware to handle ratelimits and to log requests
	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			if request.URL.Path == "/favicon.ico" {
				handler.ServeHTTP(writer, request)
				return
			}

			// check if the request is ratelimited
			if _, ok := RateLimits[request.RemoteAddr]; ok {
				if time.Since(RateLimits[request.RemoteAddr]) < time.Millisecond*time.Duration(*RateLimitTime) {
					writer.WriteHeader(http.StatusTooManyRequests)
					_, err := writer.Write([]byte(`{"error": "ratelimited"}`))
					if err != nil {
						log.Println("Error writing response: " + err.Error())
					}
					return
				}
			}

			// add ratelimit to the map
			RateLimits[request.RemoteAddr] = time.Now()

			f := &middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags), NoColor: false}
			entry := f.NewLogEntry(request)
			ww := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)
			t1 := time.Now()
			defer func() {
				entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
			}()

			handler.ServeHTTP(ww, request)
		})
	})

	err := setupEndpoints()
	if err != nil {
		log.Fatal("Error setting up API endpoints: " + err.Error())
	}

	log.Println("API endpoints initialized")
}

func main() {
	var err error
	datalocation := "data"

	// create the data folder if it doesn't exist
	// this is the location all our backend data will be stored at
	err = util.CreateFolderIfNotExists(datalocation)
	if err != nil {
		log.Fatal("Error creating data folder: " + err.Error())
	}
	log.Println("Data folder created")

	// initiate connection to the database
	// this will also create any necessary tables if they don't exist
	err = database.Startup(datalocation)
	if err != nil {
		log.Fatal("Error starting up database: " + err.Error())
	}
	log.Println("Database started")

	// finally, start listening for connections on the specified port
	log.Println("Listening for API requests at port " + *BackendPort)
	err = http.ListenAndServe(":"+*BackendPort, router)
	if err != nil {
		log.Fatal(err)
	}
}

/*
setupEndpoints ~ registers the endpoints for the api to allow response
*/
func setupEndpoints() error {
	router.Get("/favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
		// returns the favicon.ico file
		http.ServeFile(writer, request, "public/favicon.ico")
	})
	router.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		// get the content type from request
		contentType := request.Header.Get("Accept")
		// if the content type is application/json, return a json response
		if contentType == "application/json" {
			writer.WriteHeader(http.StatusNotFound)
			_, err := writer.Write([]byte(`{"error": "not found"}`))
			if err != nil {
				log.Println("Error writing response: " + err.Error())
			}
			return
		}
		// otherwise, return a html response
		writer.WriteHeader(http.StatusNotFound)
		http.ServeFile(writer, request, "public/404.html")
	})
	router.MethodNotAllowed(func(writer http.ResponseWriter, request *http.Request) {
		// get the content type from request
		contentType := request.Header.Get("Accept")
		// if the content type is application/json, return a json response
		if contentType == "application/json" {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			_, err := writer.Write([]byte(`{"error": "method not allowed"}`))
			if err != nil {
				log.Println("Error writing response: " + err.Error())
			}
			return
		}
		// otherwise, return a html response
		writer.WriteHeader(http.StatusMethodNotAllowed)
		http.ServeFile(writer, request, "public/405.html")
	})

	router.Post("/api/login", api.LoginFunction)
	router.Post("/api/register", api.RegisterFunction)
	router.Post("/api/createBook", api.CreateBookFunction)
	router.Post("/api/getBooks", api.GetBooksFunction)
	router.Post("/api/claimBook", api.ClaimBookFunction)
	router.Post("/api/unclaimBook", api.UnclaimBookFunction)

	return nil
}
