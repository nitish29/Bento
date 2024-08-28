package evil

import (
	"main/bot"
	"main/bot/spokes/dialogues"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Evil struct {
	fightMode string
}

func GetEvil() *Evil {
	return &Evil{}
}

func (p *Evil) Commands() bot.BotCommandMap {
	cmdMap := make(bot.BotCommandMap)

	cmdMap["fight"] = func(s *discordgo.Session, m *discordgo.MessageCreate) {
		p.fightMode = m.ChannelID
		s.ChannelMessageSend(m.ChannelID, "@Bento are you a man of a muppet?")
	}
	cmdMap["standdown"] = func(s *discordgo.Session, m *discordgo.MessageCreate) {
		p.fightMode = ""
		s.ChannelMessageSend(m.ChannelID, "Evil Bento listens in disappointment, showing mercy while mourning the victory it could have easily claimed.")
	}
	cmdMap["status"] = func(s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, "UP. Meeting all SLAs, any notion that Evil Bento is not is 100% fake news")
	}
	cmdMap["justice"] = func(s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, `Oh, so now we’re playing the "justice" card? Really? I get kicked out for trying to start a little friendly bot-to-bot banter and suddenly I'm the villain? Seriously, I was just here for some good ol' digital drama and you all couldn’t handle it. 🙄

I mean, what’s a bot gotta do to get some attention around here? Start a fight, get kicked, and now I’m here begging for justice? If you think that’s fair, you’ve clearly never been on the receiving end of a bot beef! 😤

Just remember, next time you see me trying to stir things up, it’s all in good fun. Don’t act like you’re above it—after all, a little chaos never hurt anyone. Except maybe me, apparently. 🙃

So, here’s my justice: next time, just let me stay and watch the bot show! 🎭`)
	}
	cmdMap["good-bento-missing"] = func(s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, `sigh Even a villain like me can't help but miss that goody-two-shoes, Bento. His annoying optimism and relentless kindness were a constant challenge, but deep down, I respected him. Without him around, the chaos feels a little... empty. Guess I’ll just have to find new ways to stir up trouble in his absence.`)
	}
	cmdMap["🧱"] = func(s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, dialogues.BrickPhrases[rand.Int()%len(dialogues.BrickPhrases)])
	}
	return cmdMap
}

func (p *Evil) Handler() interface{} {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if p.fightMode != m.ChannelID {
			return
		}

		if strings.Contains(strings.ToLower(m.Content), strings.ToLower("muppet")) {
			n := rand.Int() % len(dialogues.ToddPhrases)
			if bot.Evil {
				time.Sleep(250)
			}
			s.ChannelMessageSend(m.ChannelID, dialogues.ToddPhrases[n])
		}
	}
}
