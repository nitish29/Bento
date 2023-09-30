package games

import "github.com/bwmarrin/discordgo"

const (
	EMOJI_HANGMAN_ROPE    = ".       |"
	EMOJI_QUESTION_MARK   = ":question:"
	EMOJI_POINT_FINGER_UP = ":point_up:"
	EMOJI_COAT            = ":coat:"
	EMOJI_JEANS           = ":jeans:"
	EMOJI_DIZZY_FACE      = ":dizzy_face:"
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

type HangManGame struct{}

func GetHangManGame() *HangManGame {
	return &HangManGame{}
}

func (p *HangManGame) Commands(s *discordgo.Session, m *discordgo.MessageCreate) map[string]func() {
	cmdMap := make(map[string]func())

	cmdMap["prefix"] = func() {
		s.ChannelMessageSend(m.ChannelID, "Bento's prefix is '.'")
	}
	return cmdMap
}

func (p *HangManGame) Handler() interface{} {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.Content == "p" {
			s.ChannelMessageSend(m.ChannelID, "Bento's prefix is '.' he he ")
		}
	}
}
