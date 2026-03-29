package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"discord-rpg-bot/db"
	"discord-rpg-bot/game"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session  *discordgo.Session
	DB       *db.Database
	Handlers *Handlers
}

func New(token string, database *db.Database) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	handlers := &Handlers{
		DB:            database,
		CombatManager: game.NewCombatManager(),
	}

	b := &Bot{
		Session:  session,
		DB:       database,
		Handlers: handlers,
	}

	session.AddHandler(handlers.HandleInteraction)
	session.Identify.Intents = discordgo.IntentsGuildMessages

	return b, nil
}

func (b *Bot) Start() error {
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}

	log.Println("Registering slash commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, cmd := range commands {
		registered, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, "", cmd)
		if err != nil {
			log.Printf("Error registering command '/%s': %v", cmd.Name, err)
			continue
		}
		registeredCommands[i] = registered
		log.Printf("Registered: /%s", cmd.Name)
	}

	fmt.Println("========================================")
	fmt.Println("  Legends of Eldoria Bot is running!")
	fmt.Println("  Use /adventure in Discord to play")
	fmt.Println("  Press Ctrl+C to stop")
	fmt.Println("========================================")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	log.Println("Shutting down...")

	for _, cmd := range registeredCommands {
		if cmd != nil {
			if err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, "", cmd.ID); err != nil {
				log.Printf("Error deleting command '/%s': %v", cmd.Name, err)
			}
		}
	}

	return b.Session.Close()
}
