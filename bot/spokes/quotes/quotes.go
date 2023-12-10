package quotes

import (
	"encoding/json"
	"fmt"
	"io"
	"main/bot"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Quote struct {
	bot.DefaultSpoke
}

func GetQuote() *Quote {
	return &Quote{}
}

type QuoteResponse struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

func getRapidToken() string {
	token := os.Getenv("RAPID_API_TOKEN")
	if len(token) == 0 {
		return ""
	}
	return token
}

func (Quote) Commands(s *discordgo.Session, m *discordgo.MessageCreate) map[string]func() {
	cmdMap := make(map[string]func())

	cmdMap["quote"] = func() {
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
		s.ChannelMessageSend(m.ChannelID, msg)
	}
	return cmdMap

}
