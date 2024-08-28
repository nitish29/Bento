package main

import (
	"log/slog"
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
		slog.Error("Error creating Discord session", "err", err)
	}

	// Register spokes to bot
	if bot.Evil {
		b.RegisterSpoke(evil.GetEvil())
	} else {
		b.RegisterSpoke(dialogues.GetDialogues())
		jobs.StartJob(b.Session)
	}

	b.RegisterSpoke(general.GetPrefix())
	b.RegisterSpoke(hangman.GetHangManSpoke())

	b.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	b.SyncSpokes()

	err = b.Open()
	if err != nil {
		slog.Error("Error opening Discord session", "err", err)
	}
	defer func() {
		slog.Info("Bot terminating")
		b.Close()
	}()

	slog.Info("Bot running", "user", b.State.User)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
