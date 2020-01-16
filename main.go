package main

import (
	"github.com/cloudmusic-dev/backend/api"
	"github.com/cloudmusic-dev/backend/configuration"
	"github.com/cloudmusic-dev/backend/database"
	"github.com/goccy/go-yaml"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)



func main() {
	// Load configuration
	configurationSource, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	configuration := configuration.Configuration{}
	err = yaml.Unmarshal(configurationSource, &configuration)
	if err != nil {
		log.Fatalf("Failed to parse configuration: %v", err)
	}

	// Connect to database
	database.InitializeDatabase(configuration)
	defer database.CloseDatabase()

	// Create our main api router
	router := mux.NewRouter()
	api.CreatePlaylistRouter(router.PathPrefix("/api").Subrouter())
	api.CreateFolderRouter(router.PathPrefix("/api").Subrouter())

	server := &http.Server{
		Handler:      router,
		Addr:         ":32546",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start web server
	log.Fatal(server.ListenAndServe())
}
