package games

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
