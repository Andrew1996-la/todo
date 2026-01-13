package main

import (
	"fmt"
	"log"
	"os"
	"todo/internal/database"
	"todo/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	serverPort := os.Getenv("SERVER_PORT")

	if databaseUrl == "" || serverPort == "" {
		log.Fatal("DATABASE_URL or SERVER_PORT not set")
	}

	log.Printf("Starting server on port %s", serverPort)

	db, err := database.Connect(databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	taskStore := database.NewTaskStore(db)
	handlersTodo := handlers.NewTodoHandler(taskStore)
	server := handlers.NewHTTPServer(handlersTodo, serverPort)
	server.Start()

	defer func() {
		err := db.Close()
		if err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	fmt.Println("Connected to database")
}
