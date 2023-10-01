package games

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/slices"
)

const (
	EMOJI_HANGMAN_ROPE    = ".       |"
	EMOJI_QUESTION_MARK   = ":question:"
	EMOJI_POINT_FINGER_UP = ":point_up:"
	EMOJI_COAT            = ":coat:"
	EMOJI_JEANS           = ":jeans:"
	EMOJI_DIZZY_FACE      = ":dizzy_face:"
	STR_GAME_OVER         = "\n:regional_indicator_g: :regional_indicator_a: :regional_indicator_m: :regional_indicator_e:      :regional_indicator_o: :regional_indicator_v: :regional_indicator_e: :regional_indicator_r:\n "
	STR_YOU_WON           = "\n:regional_indicator_y: :regional_indicator_o: :regional_indicator_u:     :regional_indicator_w: :regional_indicator_o: :regional_indicator_n:\n"
	STR_YOU_LOST          = "\n:regional_indicator_y: :regional_indicator_o: :regional_indicator_u:     :regional_indicator_l: :regional_indicator_o: :regional_indicator_s: :regional_indicator_t:" + "\n"
	STR_GET_WORD          = "Tell me the secret word, type \".word [your word]\" "
	STR_ALREADY_PLAYING   = "I am already playing a game"
	STR_ABORT_SUCCESS     = "Game stopped"
	STR_NO_GAME           = "No game running"
	STR_GUESSED_SO_FAR    = "Guessed so far : "
	STR_TRY_HANGMAN       = "Please use .hangman before using the 'word' command"
	STR_WORD_EMPTY        = "Please type your word with the command Eg: .word [your word]"
)

type Hangman struct {
	serverId           string
	channeld           string
	challenger         string
	isAcceptingLetters bool
	isGameOver         bool
	hasPlayerWon       bool
	game               *Game
}

type Game struct {
	truth          string
	guessedLetters []string
	wrongAnswers   int
	blanks         []string
	status         string
	printList      []string
}

func New(serverId string, channelId string, challenger string) *Hangman {
	return &Hangman{
		serverId:           serverId,
		channeld:           channelId,
		challenger:         challenger,
		isAcceptingLetters: false,
		isGameOver:         false,
		hasPlayerWon:       false,
		game:               initGame(),
	}
}

func initGame() *Game {
	return &Game{
		truth:          "",
		guessedLetters: make([]string, 0),
		wrongAnswers:   0,
		blanks:         make([]string, 0),
		status:         "",
		printList: []string{EMOJI_HANGMAN_ROPE, ".     " + EMOJI_DIZZY_FACE,
			"." + EMOJI_POINT_FINGER_UP, EMOJI_COAT, EMOJI_POINT_FINGER_UP,
			" .     " + EMOJI_JEANS},
	}
}

func (h *Hangman) updateHangmanStatus(wrongAnswers int) {
	// The arms and body needs to be in a straight line and therefore no new line character
	if wrongAnswers == 3 || wrongAnswers == 4 {
		h.game.status += h.game.printList[wrongAnswers-1]
	} else {
		h.game.status += h.game.printList[wrongAnswers-1] + "\n"
	}
	if h.game.wrongAnswers == len(h.game.printList) {
		h.isGameOver = true
		h.hasPlayerWon = false
	}
}

func (h *Hangman) updateBlanks(letter string) {
	if letter == "" {
		for _, _ = range h.game.truth {
			h.game.blanks = append(h.game.blanks, EMOJI_QUESTION_MARK)
		}
	} else {
		for i, v := range h.game.truth {
			if letter == string(v) {
				h.game.blanks[i] = string(v)
			}
		}
	}

	if strings.Join(h.game.blanks, "") == h.game.truth {
		h.isGameOver = true
		h.hasPlayerWon = true
	}
}

func (h *Hangman) processInput(letter string) {
	if h.isGameOver {
		return
	}

	if strings.Contains(h.game.truth, letter) {
		h.updateBlanks(letter)
	} else {
		h.game.wrongAnswers += 1
		h.updateHangmanStatus(h.game.wrongAnswers)
	}
	if !slices.Contains(h.game.guessedLetters, letter) && letter != "" {
		h.game.guessedLetters = append(h.game.guessedLetters, letter)
	}

}

func (h *Hangman) getGameStatus() string {
	result := ""
	if h.isGameOver {
		result += STR_GAME_OVER
		if h.hasPlayerWon {
			result += STR_YOU_WON
		} else {
			result += STR_YOU_LOST + "\n" + "The Word was " + h.game.truth
		}
	} else {
		result = STR_GUESSED_SO_FAR
		guesses := h.game.guessedLetters
		for _, v := range guesses {
			result += v + ","
		}
	}

	blanks := ""
	for _, b := range h.game.blanks {
		blanks += b + "\t"
	}

	return "\n" + blanks + "\n" + h.game.status + "\n" + result
}

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

	cmdMap["hangman"] = func() {
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

func (p *HangManSpoke) wordCmd(s *discordgo.Session, m *discordgo.MessageCreate) func() {
	return func() {
		author := m.Author.ID
		challenger, _ := s.UserChannelCreate(author)
		// TODO : handle case where user starts game in 2 servers. this might get confused in that case.

		serverID, exists := p.challengerToServer[challenger.ID]
		if !exists {
			s.ChannelMessageSend(m.ChannelID, STR_TRY_HANGMAN)
			return
		}
		gameInstance := p.gameInstances[serverID]

		content := strings.Split(m.Content, " ")

		if len(content) <= 1 {
			s.ChannelMessageSend(m.ChannelID, STR_WORD_EMPTY)
			return
		}
		gameInstance.isAcceptingLetters = true
		gameInstance.game.truth = content[1]
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
