package main

import (
	"fmt"
	"main/bot"

	// "main/bot/jobs"
	"main/bot/jobs"
	"main/bot/spokes/dialogues"
	"main/bot/spokes/evil"
	"main/bot/spokes/games/hangman"
	"main/bot/spokes/general"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func main() {
	b, err := bot.New()
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
	}

	// Register spokes to bot

	if bot.Evil {
		b.RegisterSpoke(evil.GetEvil())
	} else {
		b.RegisterSpoke(dialogues.GetDialogues())
	}
	b.RegisterSpoke(general.GetPrefix())
	b.RegisterSpoke(hangman.GetHangManSpoke())

	b.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	b.SyncSpokes()

	err = b.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}
	defer func() {
		fmt.Println("Bot terminating")
		b.Close()
	}()

	if !bot.Evil {
		jobs.StartJob(b.Session)
	}

	fmt.Println("Bot running")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
