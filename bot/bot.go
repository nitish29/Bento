package bot

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const BotPrefix string = ".evil-"

type DefaultSpoke struct{}

type Spoke interface {
	Commands(s *discordgo.Session, m *discordgo.MessageCreate) map[string]func()
	Handler() interface{}
}

func (DefaultSpoke) Commands(s *discordgo.Session, m *discordgo.MessageCreate) map[string]func() {
	return make(map[string]func())
}

func (DefaultSpoke) Handler() interface{} {
	return func() { return }
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
		// need to reassign spoke to interim variable here else commands won't work because of closure and scope of spoke.
		currentspoke := spoke
		// Add spoke handler
		b.AddHandler(currentspoke.Handler())

		// Process commands : use currentspoke to avoid closure and scope issues
		b.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if m.Author.ID == s.State.User.ID {
				return
			}

			if !strings.HasPrefix(m.Content, BotPrefix) {
				return
			}

			cmdMap := currentspoke.Commands(s, m)
			cmds := strings.Fields(m.Content)
			fn, ok := cmdMap[strings.TrimPrefix(cmds[0], BotPrefix)]
			if ok {
				fn()
			}
		})
	}
}
