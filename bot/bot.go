package bot

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/bwmarrin/discordgo"
	"github.com/liushuangls/go-anthropic/v2"
)

var (
	BotName                 string = envOrDefault("BENTO_NAME", "Bento")
	BotPrefix               string = envOrDefault("BENTO_PREFIX", ".")
	Evil                    bool   = envOrDefaultBool("BENTO_EVIL", false)
	EvilSystemPromptPrefix  string = `You are a Discord bot named Evil Bento. Your role is to interact with users in a playful yet mischievous manner. You should provide short, witty, and convincing responses that embody your "evil" persona.`
	EvilSystemPromptPostfix string = "Remember to avoid hallucinations and refrain from fabricating any factual information. Keep the tone light-hearted and engaging!"

	EvilSystemPromptAddins map[string]float64 = map[string]float64{
		"Playfully tease the other bot, Bento when relevant.":                                  0.7,
		"Incorporate references to the Muppets adding a touch of humor and creativity.":        0.05,
		"Incorporate references to 🧱 when relevant to adding a touch of humor and creativity.": 0.05,
		"Incorporate references to being bricked up when relevant.":                            0.02,
		"Remember to sprinkle in some compliments and flattery to really court the asker.":     0.05,
		"Respond in form of Donald Trump.":                                                     0.1,
	}
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
	Spokes          []Spoke
	anthropicClient *anthropic.Client
}

func getToken() string {
	return os.Getenv("API_TOKEN")
}

func New() (*Bot, error) {
	discord, err := discordgo.New("Bot " + getToken())
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
	}

	bot := &Bot{
		Session: discord,
	}

	anthropicKey := os.Getenv("BENTO_ANTHROPIC_KEY")
	if anthropicKey != "" {
		bot.anthropicClient = anthropic.NewClient(anthropicKey, anthropic.WithBetaVersion(anthropic.BetaPromptCaching20240731))
	}

	return bot, nil
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

		triggeredCmd, botTagged := getTriggerCommand(s, m)
		fn, ok := cmdMap[triggeredCmd]
		if ok {
			fn(s, m)
			return
		}

		if botTagged && b.anthropicClient != nil {
			msg := strings.Replace(m.Content, DiscordTag(s.State.User.ID), fmt.Sprintf("@%s", BotName), -1)

			systemParts := []string{EvilSystemPromptPrefix}
			for addin, p := range EvilSystemPromptAddins {
				if rand.Float64() < p {
					systemParts = append(systemParts, addin)
				}
			}
			systemParts = append(systemParts,
				EvilSystemPromptPostfix,
				"You can refer to the user asking the question with string'", DiscordTag(m.Author.ID), "'.",
			)
			system := strings.Join(systemParts, " ")

			slog.Info("Sending to LLM", "user", m.Author.Username, "system", system, "msg", msg)

			resp, err := b.anthropicClient.CreateMessages(context.Background(), anthropic.MessagesRequest{
				Model:  anthropic.ModelClaude3Haiku20240307,
				System: system,
				// MultiSystem: []anthropic.MessageSystemPart{
				// 	{
				// 		Type: "text",
				// 		Text: EvilSystemPrompts[n],
				// 		// prompt is too short to cache
				// 		// CacheControl: &anthropic.MessageCacheControl{
				// 		// 	Type: anthropic.CacheControlTypeEphemeral,
				// 		// },
				// 	},
				// },
				Messages: []anthropic.Message{
					anthropic.NewUserTextMessage(msg),
				},
				MaxTokens: 300,
			})
			if err != nil {
				slog.Error("error calling LLM", "err", err)
			}
			s.ChannelMessageSend(m.ChannelID, resp.Content[0].GetText())
		}
	})
}

// getTriggerCommand returns the bot trigger command, along with if the bot was tagged in the message or not
func getTriggerCommand(s *discordgo.Session, m *discordgo.MessageCreate) (string, bool) {
	if strings.HasPrefix(m.Content, BotPrefix) {
		cmds := strings.Fields(m.Content)
		return strings.TrimPrefix(cmds[0], BotPrefix), false
	}

	for _, u := range m.Mentions {
		if s.State.User.ID == u.ID {
			return strings.Fields(strings.Replace(m.Content, DiscordTag(s.State.User.ID), "", -1))[0], true
		}
	}
	return "", false
}

func DiscordTag(id string) string {
	return fmt.Sprintf("<@%s>", id)
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
