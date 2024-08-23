package general

import (
	"github.com/bwmarrin/discordgo"
)

type Prefix struct{}

func GetPrefix() *Prefix {
	return &Prefix{}
}

func (p *Prefix) Commands(s *discordgo.Session, m *discordgo.MessageCreate) map[string]func() {
	cmdMap := make(map[string]func())

	cmdMap["prefix"] = func() {
		s.ChannelMessageSend(m.ChannelID, "Evil Bento's prefix is '.evil-'")
	}
	return cmdMap
}

func (p *Prefix) Handler() interface{} {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.Content == "p" {
			s.ChannelMessageSend(m.ChannelID, "Evil Bento's prefix is '.evil-' he he ")
		}
	}
}
