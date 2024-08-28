package general

import (
	"fmt"
	"main/bot"

	"github.com/bwmarrin/discordgo"
)

type Prefix struct{}

func GetPrefix() *Prefix {
	return &Prefix{}
}

func (p *Prefix) Commands(s *discordgo.Session, m *discordgo.MessageCreate) map[string]func() {
	cmdMap := make(map[string]func())

	cmdMap["prefix"] = func() {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s's prefix is '%s'", bot.BotName, bot.BotPrefix))
	}
	return cmdMap
}

func (p *Prefix) Handler() interface{} {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.Content == "p" {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s's prefix is '%s' he he ", bot.BotName, bot.BotPrefix))
		}
	}
}
