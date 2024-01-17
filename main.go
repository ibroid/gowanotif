package main

import (
	"gowhatsapp/database"
	"gowhatsapp/events"
	"gowhatsapp/log"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	log.LogInit()
	database.InitDBLocal()

	events.StartSippEvent()
}
