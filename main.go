package main

import (
	"fmt"
	"main/bot"
	"main/bot/spokes/dialogues"
	"main/bot/spokes/general"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func main() {
	bot, err := bot.New()
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
	}

	// Register spokes to bot

	bot.RegisterSpoke(dialogues.GetDialogues())
	bot.RegisterSpoke(general.GetPrefix())

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
