package general

import "github.com/bwmarrin/discordgo"

type Prefix struct{}

func GetPrefix() *Prefix {
	return &Prefix{}
}

func (p *Prefix) Commands() map[string]interface{} {
	return nil
}

func (p *Prefix) Subcommands() map[string]interface{} {
	return nil
}

func (p *Prefix) Description() string {
	return ""
}

func (p *Prefix) Handler() interface{} {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == ".prefix" {
			s.ChannelMessageSend(m.ChannelID, "Bento's prefix is .")
		}
	}

}
