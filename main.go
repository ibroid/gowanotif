package main

import (
	"gowhatsapp/database"
	"gowhatsapp/log"
	wa "gowhatsapp/whatsapp"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	log.LogInit()
	database.InitDBLocal()

	wa.InitializeThirdParty()
}
