package jobs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type QuoteResponse struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

var fedxLobby = "1124585337855950912"

var subscribedChannels = []string{"207372726221012993", "399577277903536138", fedxLobby}

func getRapidToken() string {
	token := os.Getenv("RAPID_API_TOKEN")
	if len(token) == 0 {
		return ""
	}
	return token
}

func StartJob(s *discordgo.Session) {
	QueryQuote(s)
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				QueryQuote(s)
			}
		}
	}()
}

func QueryQuote(s *discordgo.Session) {
	rapidAPIKey := getRapidToken()
	rapidAPIHost := "quotes-inspirational-quotes-motivational-quotes.p.rapidapi.com"

	url := "https://quotes-inspirational-quotes-motivational-quotes.p.rapidapi.com/quote"

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	query := req.URL.Query()
	query.Add("token", "ipworld.info")
	req.URL.RawQuery = query.Encode()

	req.Header.Set("X-RapidAPI-Key", rapidAPIKey)
	req.Header.Set("X-RapidAPI-Host", rapidAPIHost)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var quoteResponse QuoteResponse
	err = json.Unmarshal(body, &quoteResponse)
	if err != nil {
		return
	}

	msg := fmt.Sprintf(" :bulb: Quote of the day :bulb:\n > %s", quoteResponse.Text)


	for _, c := range subscribedChannels {
		if c == fedxLobby {
			msg += "\n\n :tada: Happy Birthday, Justin! :tada:"
		}
		s.ChannelMessageSend(c, msg)
	}
}
