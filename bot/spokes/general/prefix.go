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

func (p *Prefix) Commands() bot.BotCommandMap {
	cmdMap := make(bot.BotCommandMap)

	cmdMap["prefix"] = func(s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s's prefix is '%s'", bot.BotName, bot.BotPrefix))
	}
	return cmdMap
}

func (p *Prefix) Handlers() []interface{} {
	return []interface{}{
		func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if m.Author.ID == s.State.User.ID {
				return
			}
			if m.Content == "p" {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s's prefix is '%s' he he ", bot.BotName, bot.BotPrefix))
			}
		},
	}
}
