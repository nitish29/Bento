package hangman

import (
	"strings"

	"golang.org/x/exp/slices"
)

const (
	EMOJI_HANGMAN_ROPE    = "\n.       |"
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
	STR_GUESSED_SO_FAR    = "Guessed so far : \n"
	STR_TRY_HANGMAN       = "Please use .hangman before using the 'word' command"
	STR_WORD_EMPTY        = "Please type your word with the command Eg: .word [your word]"
	STR_WORD_EXISTS       = "I already have a word"
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
		for range h.game.truth {
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
	letter = strings.ToUpper(letter)
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
