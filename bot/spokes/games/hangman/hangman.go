package hangman

import (
	"strings"

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

func (h *HangManSpoke) Commands(s *discordgo.Session, m *discordgo.MessageCreate) map[string]func() {
	cmdMap := make(map[string]func())
	cmdMap["hangman"] = h.hangmanCmd(s, m)
	cmdMap["word"] = h.wordCmd(s, m)
	cmdMap["abort"] = h.abortCmd(s, m)
	return cmdMap
}

func (h *HangManSpoke) Handler() interface{} {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
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
}

func (h *HangManSpoke) wordCmd(s *discordgo.Session, m *discordgo.MessageCreate) func() {
	return func() {
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
		gameInstance.isAcceptingLetters = true
		gameInstance.game.truth = strings.ToUpper(content[1])
		gameInstance.processInput("")
		// Send blanks and start the game
		s.ChannelMessageSend(gameInstance.channeld, gameInstance.getGameStatus())
	}
}

func (h *HangManSpoke) abortCmd(s *discordgo.Session, m *discordgo.MessageCreate) func() {
	return func() {
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
}

func (h *HangManSpoke) hangmanCmd(s *discordgo.Session, m *discordgo.MessageCreate) func() {
	return func() {
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
}
