package dialogues

import (
	"main/bot"
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
	"Muppets are taking over the world!",
	"I’m living in a Muppet extravaganza.",
	"Everywhere I look, it’s Muppet mayhem.",
	"Muppets, Muppets everywhere!",
	"Life’s better with a Muppet twist.",
	"It’s a Muppet-palooza!",
	"Muppets on the brain, can't escape them!",
	"Welcome to the Muppet madness.",
	"It’s a Muppet kind of day.",
	"Muppet mania is in full swing.",
	"Can’t get enough of these Muppets!",
	"Muppets galore in every corner.",
	"A Muppet for every occasion!",
	"Just another day in Muppetland.",
	"Muppet craziness is all around.",
	"Drowning in a sea of Muppets!",
	"This place is Muppet central.",
	"Feeling like a Muppet today.",
	"Muppets are everywhere you turn.",
	"It’s a Muppet party!",
	"Caught in a Muppet whirlwind.",
	"All Muppets, all the time.",
	"Muppets have taken over the scene.",
	"Everything’s better with Muppets.",
	"Muppets make everything more fun.",
	"Lost in a world of Muppets.",
	"Muppet takeover in progress.",
	"The Muppet invasion is real.",
	"Muppets galore and loving it!",
	"Life in Muppetville is wild.",
	"Surrounded by lovable Muppets.",
	"Every day’s a Muppet adventure.",
	"It’s Muppet mania out here!",
	"The Muppet magic is in full effect.",
	"Welcome to the land of Muppets.",
	"It’s a Muppet-filled wonderland.",
	"Muppets make the world go round.",
	"Can’t escape the Muppet charm.",
	"Living in a Muppet dreamland.",
	"The Muppets have taken over.",
	"Endless Muppet fun and frolic.",
	"Muppet madness is everywhere!",
	"Muppet overload—what a thrill!",
	"Everywhere you go, Muppets follow.",
	"Muppets are taking center stage.",
	"In a world of Muppet whimsy.",
	"Muppets are my daily dose of joy.",
	"Feeling the Muppet vibe all around.",
	"Muppet extravaganza at every turn.",
	"Muppets are always in style.",
	"A Muppet-filled day is a good day.",
	"Lost in Muppet wonderland.",
	"Muppets bringing smiles everywhere.",
	"Muppets are in the spotlight today.",
	"It’s a Muppet kind of party.",
	"Caught in a Muppet whirlwind.",
	"The Muppet parade has begun!",
	"Everywhere you look, Muppets!",
	"Muppet magic is in the air.",
	"The Muppets are out and about.",
	"Muppet delight is all around.",
	"It’s Muppet overload, and I love it!",
	"Muppet joy in every nook and cranny.",
	"Living in a Muppet dream.",
	"Surrounded by Muppet cheer.",
	"Muppets have taken over the scene.",
	"Can’t escape the Muppet vibe.",
	"Muppets everywhere, and it’s fantastic.",
	"It’s a Muppet world, and we’re living in it.",
	"Feeling the Muppet magic all around.",
	"Every day’s a Muppet celebration.",
	"Muppet fun is in the air.",
	"Caught in a whirlwind of Muppets.",
	"Muppets make life more entertaining.",
	"Inundated with Muppet joy.",
	"The Muppets have arrived in style.",
	"Muppet fun at every corner.",
	"Embracing the Muppet invasion!",
	"The Muppets are everywhere you look.",
	"Living in a Muppet fairy tale.",
	"Muppet fun is a constant companion.",
	"Muppets turning every day into a party.",
	"It’s a Muppet celebration all the time.",
	"Muppets are the stars of the show.",
	"Everywhere you turn, there’s a Muppet.",
	"Muppet magic is all around us.",
	"Surrounded by Muppet cheerfulness.",
	"The Muppets make everything brighter.",
	"In the midst of Muppet mania.",
	"Muppet wonderland, here we come!",
	"Muppets are the life of the party.",
	"Can’t get away from the Muppet fun.",
	"Living in a Muppet extravaganza.",
	"Every day is a Muppet adventure.",
	"The Muppets are in full force.",
	"Muppet joy is everywhere.",
	"Muppet fun knows no bounds.",
	"Life with Muppets is never dull.",
	"A world filled with Muppet magic.",
	"Muppets make every moment special.",
}
var Bento = "My creator named me after Ben(ben) and Todd(to), two great minds. One is scary clever and the other is cleverly funny"

type Dialogues struct{}

func GetDialogues() *Dialogues {
	return &Dialogues{}
}

func (p *Dialogues) Commands() bot.BotCommandMap {
	cmdMap := make(bot.BotCommandMap)

	cmdMap["dialogues"] = func(s *discordgo.Session, m *discordgo.MessageCreate) {
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
