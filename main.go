package main

import (
	"fmt"
	"main/bot"
	"main/bot/dialogues"
	"main/bot/spokes/general"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const botPrefix string = "."

func main() {
	bot, err := bot.New()
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
	}

	bot.RegisterSpoke(general.GetPrefix())

	bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		commmands := strings.Split(m.Content, " ")
		baseCmd := commmands[0]
		if len(commmands) > 1 {
			if baseCmd == botPrefix {
				switch commmands[1] {
				case "name":
					s.ChannelMessageSend(m.ChannelID, dialogues.Bento)
				}
			}
		}

		if strings.Contains(strings.ToLower(m.Content), strings.ToLower("muppet")) {
			rand.Seed(time.Now().Unix())
			n := rand.Int() % len(dialogues.ToddPhrases)
			s.ChannelMessageSend(m.ChannelID, dialogues.ToddPhrases[n])
		}

		if strings.Contains(strings.ToLower(m.Content), strings.ToLower("oops")) {
			rand.Seed(time.Now().Unix())
			n := rand.Int() % len(dialogues.BenPhrases)
			s.ChannelMessageSend(m.ChannelID, dialogues.BenPhrases[n])
		}
	})

	bot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	bot.SyncSpokes()

	err = bot.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}
	defer func() {
		fmt.Println("Bot terminating")
		bot.Close()
	}()

	fmt.Println("Bot running")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
