package bot

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Spoke interface {
	Commands() map[string]interface{}
	Subcommands() map[string]interface{}
	Description() string
	Handler() interface{}
}

type Bot struct {
	*discordgo.Session
	Spokes []Spoke
}

func getToken() string {
	token := os.Getenv("API_TOKEN")
	if len(token) == 0 {
		return ""
	}
	return token
}

func New() (*Bot, error) {
	discord, err := discordgo.New("Bot " + getToken())
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
	}
	return &Bot{
		Session: discord,
	}, nil
}

func (b *Bot) RegisterSpoke(spoke Spoke) {
	b.Spokes = append(b.Spokes, spoke)
}

func (b *Bot) SyncSpokes() {
	for _, spoke := range b.Spokes {
		// Add spoke handler
		b.AddHandler(spoke.Handler())
		// Process commands
	}
}
