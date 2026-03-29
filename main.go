package main

import (
	"log"
	"os"

	"discord-rpg-bot/bot"
	"discord-rpg-bot/db"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}

	database, err := db.NewDatabase("eldoria.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	b, err := bot.New(token, database)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if err := b.Start(); err != nil {
		log.Fatalf("Bot error: %v", err)
	}
}
