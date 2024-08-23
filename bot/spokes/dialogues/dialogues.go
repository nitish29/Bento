package dialogues

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var BenPhrases = []string{"Not ideal", "yo oh!"}
var ToddPhrases = []string{
	"It's like muppets in space", 
	"Surprise everyone is a muppet", 
	"I'm surrounded by Muppets", 
	"Muppets to the right, Muppets to the left", 
	"MAN or MUPPET?",
	"When did I join the muppet show?",
	"The muppet uprising has begun",
    "Muppets are the new normal",
    "It's muppets all the way down",
    "Too many muppets in the kitchen",
    "Muppets, muppets everywhere",
    "Muppets in disguise",
    "The muppet parade never ends",
    "When did I join the muppet show?",
    "A muppet by any other name",
    "The land of muppets and mayhem",
    "Muppets on a mission",
    "Muppets, assemble!",
    "Is it a man or a muppet?",
    "Muppet logic is the only logic",
    "Lost in a sea of muppets",
    "Living in a muppet world",
    "Muppets taking over the world",
    "Muppets as far as the eye can see",
    "Muppet central calling",
    "When in doubt, blame the muppets",
    "Welcome to the muppet madness",
    "The muppets are strong with this one",
    "Muppets on the loose",
    "Unleash the muppets",
    "Muppet chaos incoming",
    "This is a muppet production",
    "Muppet mayhem in progress",
    "Muppets on parade",
    "Muppets gone wild",
    "When muppets attack",
    "The muppet conspiracy",
    "Behind every corner, a muppet",
    "It's a muppet extravaganza",
    "Muppet takeover imminent",
    "Muppet madness unleashed",
    "March of the muppets",
    "Muppets of the round table",
    "The secret life of muppets",
    "Muppet mode activated",
    "Living la vida muppet",
	}

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
