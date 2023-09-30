package dialogues

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var BenPhrases = []string{"Not ideal", "yo oh!"}
var ToddPhrases = []string{"It's like muppets in space", "Surprise everyone is a muppet", "I'm surrounded by Muppets", "Muppets to the right, Muppets to the left"}
var Bento = "My creator named me after Ben(ben) and Todd(to), two great minds. One is scary clever and the other is cleverly funny"

type Dialogues struct{}

func GetDialogues() *Dialogues {
	return &Dialogues{}
}

func (p *Dialogues) Commands(s *discordgo.Session, m *discordgo.MessageCreate) map[string]func() {
	cmdMap := make(map[string]func())

	cmdMap["dialogues"] = func() {
		s.ChannelMessageSend(m.ChannelID, "Dialogues command '.'")
	}
	return cmdMap
}

func (p *Dialogues) Handler() interface{} {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.Contains(strings.ToLower(m.Content), strings.ToLower("muppet")) {
			rand.Seed(time.Now().Unix())
			n := rand.Int() % len(ToddPhrases)
			s.ChannelMessageSend(m.ChannelID, ToddPhrases[n])
		}

		if strings.Contains(strings.ToLower(m.Content), strings.ToLower("oops")) {
			rand.Seed(time.Now().Unix())
			n := rand.Int() % len(BenPhrases)
			s.ChannelMessageSend(m.ChannelID, BenPhrases[n])
		}
	}

}
