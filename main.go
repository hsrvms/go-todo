package main

import (
	"log"
	"os"

	"github.com/hsrvms/todoapp/database"
	"github.com/hsrvms/todoapp/server"
	"github.com/hsrvms/todoapp/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURI := os.Getenv("DB_URI")
	pgStorage := database.NewPgStorage(dbURI)
	db, err := pgStorage.Init()
	if err != nil {
		log.Fatal(err)
	}

	repository := store.NewRepository(db)
	api := server.NewAPIServer(":8080", repository)
	api.Start()

}
