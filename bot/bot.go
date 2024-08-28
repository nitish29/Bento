package bot

import (
	"fmt"
	"golang.org/x/exp/maps"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	BotName   string = envOrDefault("BENTO_NAME", "Bento")
	BotPrefix string = envOrDefault("BENTO_PREFIX", ".")
	Evil      bool   = envOrDefaultBool("BENTO_EVIL", false)
)

func envOrDefaultBool(key string, defaultVal bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		panic(fmt.Sprintf("failed to parse bool from %s env %v: %v", v, key, err))
	}
	return b
}

func envOrDefault(key string, defaultVal string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	return v
}

type DefaultSpoke struct{}

type BotCommandMap = map[string]func(s *discordgo.Session, m *discordgo.MessageCreate)

type Spoke interface {
	Commands() BotCommandMap
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
	cmdMap := make(BotCommandMap)

	for _, spoke := range b.Spokes {
		// Add spoke handler
		b.AddHandler(spoke.Handler())

		m := spoke.Commands()
		for cmd, f := range m {
			cmdMap[cmd] = f
		}
	}

	cmdMap["help"] = helpResponse(maps.Keys(cmdMap))

	b.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if !strings.HasPrefix(m.Content, BotPrefix) {
			return
		}

		cmds := strings.Fields(m.Content)
		triggeredCmd := strings.TrimPrefix(cmds[0], BotPrefix)
		fn, ok := cmdMap[triggeredCmd]
		if ok {
			fn(s, m)
		}
	})
}

func helpResponse(cmdList []string) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	for i, v := range cmdList {
		cmdList[i] = "- " + BotPrefix + v
	}

	cmdString := strings.Join(cmdList, "\n")
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s commands:\n%s", BotName, cmdString))
	}
}
