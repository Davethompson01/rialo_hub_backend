package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	database "github.com/Davethompson01/rialo_hub_backend/Database"
	route "github.com/Davethompson01/rialo_hub_backend/Route"
	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Server successfully started and listening on port :8080...")
	err := godotenv.Load(".env")

	//load port
	port := os.Getenv("PORT")
	if err != nil {
		log.Fatal(err)
	}
	if port == "" {
		log.Fatalf("Failed to load port number %v", port)
	}

	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("Failed to load Database URL")
	}

	db, err := database.DatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to Load Database connection %v", err)

	}
	config := config.ApiConfig{
		DB: db,
	}

	defer config.DB.Close()

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-API-Key",
		},
		ExposedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	routes := chi.NewRouter()
	route.CreateUser(routes, &config)
	route.Task(routes, &config)
	router.Mount("/v1", routes)
	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server is already running %v", port)
	server_err := server.ListenAndServe()
	if server_err != nil {
		log.Fatal(server_err)
	}
}
