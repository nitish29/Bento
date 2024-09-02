package hangman

import (
	"fmt"
	"main/bot"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
)

type HangManSpoke struct {
	gameInstances      map[string]*Hangman
	challengerToServer map[string]string
	serverToChallenger map[string]string
}

func GetHangManSpoke() *HangManSpoke {
	return &HangManSpoke{
		gameInstances:      make(map[string]*Hangman),
		challengerToServer: make(map[string]string),
		serverToChallenger: make(map[string]string),
	}
}

func (h *HangManSpoke) Commands() bot.BotCommandMap {
	cmdMap := make(bot.BotCommandMap)
	if bot.Evil {
		cmdMap["hangman"] = func(s *discordgo.Session, m *discordgo.MessageCreate) {
			s.ChannelMessageSend(m.ChannelID, "Evil bento doesn't play games, it is a very serious bot")
		}
		cmdMap["hangman-word-domination"] = h.hangmanCmd
	} else {
		cmdMap["hangman"] = h.hangmanCmd
	}
	cmdMap["word"] = h.wordCmd
	cmdMap["abort"] = h.abortCmd
	return cmdMap
}

func (h *HangManSpoke) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" {
		return
	}
	if m.Author.ID == s.State.User.ID {
		return
	}

	channelId := m.ChannelID
	serverId := m.GuildID

	gameInstance, exists := h.gameInstances[serverId]
	if !exists {
		return
	}

	//ignore message from other channels
	if gameInstance.channeld != channelId {
		return
	}
	if !gameInstance.isAcceptingLetters {
		return
	}
	// accept single characters only
	if len(m.Content) != 1 {
		return
	}

	gameInstance.isAcceptingLetters = false
	gameInstance.processInput(m.Content)
	s.ChannelMessageSend(m.ChannelID, gameInstance.getGameStatus())
	gameInstance.isAcceptingLetters = true
	if gameInstance.isGameOver {
		// Clean up
		delete(h.challengerToServer, gameInstance.challenger)
		delete(h.serverToChallenger, serverId)
		delete(h.gameInstances, serverId)
	}
}

func (h *HangManSpoke) wordCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	author := m.Author.ID
	challenger, _ := s.UserChannelCreate(author)
	// TODO : handle case where user starts game in 2 servers. this might get confused in that case.

	serverID, exists := h.challengerToServer[challenger.ID]
	if !exists {
		s.ChannelMessageSend(m.ChannelID, STR_TRY_HANGMAN)
		return
	}
	gameInstance := h.gameInstances[serverID]

	if len(gameInstance.game.truth) > 0 {
		s.ChannelMessageSend(m.ChannelID, STR_WORD_EXISTS)
		return
	}

	content := strings.Split(m.Content, " ")

	if len(content) <= 1 {
		s.ChannelMessageSend(m.ChannelID, STR_WORD_EMPTY)
		return
	}

	if len(content) > 2 {
		s.ChannelMessageSend(m.ChannelID, "Only one word is supported currently. Please feature multi word request with the author")
		return
	}

	for _, r := range content[1] {
		fmt.Println("rune is", r)
		if !unicode.IsLetter(r) {
			s.ChannelMessageSend(m.ChannelID, "Word supports only alphabets, please try again.")
			return
		}
	}

	gameInstance.isAcceptingLetters = true
	gameInstance.game.truth = strings.ToUpper(content[1])
	gameInstance.processInput("")
	// Send blanks and start the game
	s.ChannelMessageSend(gameInstance.channeld, gameInstance.getGameStatus())
}

func (h *HangManSpoke) abortCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" {
		return
	}

	serverId := m.GuildID
	gameInstance, exists := h.gameInstances[serverId]
	if !exists {
		s.ChannelMessageSend(m.ChannelID, STR_NO_GAME)
	}

	delete(h.challengerToServer, gameInstance.challenger)
	delete(h.serverToChallenger, serverId)
	delete(h.gameInstances, serverId)
}

func (h *HangManSpoke) hangmanCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	serverId := m.GuildID
	channelId := m.ChannelID
	authorId := m.Author.ID

	// Only one instance of game running per server
	_, exists := h.gameInstances[serverId]
	if exists {
		s.ChannelMessageSend(channelId, STR_ALREADY_PLAYING)
		return
	}
	h.gameInstances[serverId] = New(serverId, channelId, authorId)
	challenger, _ := s.UserChannelCreate(authorId)
	h.challengerToServer[challenger.ID] = serverId
	h.serverToChallenger[serverId] = challenger.ID
	if challenger.ID == channelId {
		s.ChannelMessageSend(m.ChannelID, "You cannot play Hangman in a private chat.")
		return
	}
	s.ChannelMessageSend(challenger.ID, STR_GET_WORD)
}
