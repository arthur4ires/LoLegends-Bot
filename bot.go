package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var summonerReturn map[string]invocadorStruct

const (
	riotApiKey = ""
	botApiKey  = ""
	nickUrl    = "https://br.api.pvp.net/api/lol/br/v1.4/summoner/by-name/"
)

type invocadorStruct struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	SummonerLevel int    `json:"summonerLevel"`
	RevisionDate  int64  `json:"revisionDate"`
}

func htmlDownload(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return string(body)
}

func comandosBot(text string) string {
	if string(text[0]) == "/" {
		if strings.Contains(text, "/summoner") {
			userName := strings.Replace(text, "/summoner", "", -1)
			userName = strings.Replace(userName, " ", "", -1)

			response := htmlDownload(nickUrl + userName + "?api_key=" + riotApiKey)

			err := json.Unmarshal([]byte(response), &summonerReturn)

			if err != nil {
				panic(err)
			}

			if summonerReturn[userName].ID == 0 {
				return "<b>Error:</b>Summoner not found!"
			}

			return fmt.Sprintf("<b>Summoner statistics:</b>\n\n<b>Name:</b>%s \n<b>ID:</b>%d \n<b>Level:</b> %d", summonerReturn[userName].Name, summonerReturn[userName].ID, summonerReturn[userName].SummonerLevel)
		}
		if strings.Contains(text, "/start") {
			return "Hello summoner , use the command /summoner to check your statistics!"
		}
	} else {
		return "Any doubt ? Send me a /start!"
	}
	return ""
}

func main() {
	bot, err := tgbotapi.NewBotAPI(botApiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Bot autorized on acc :  %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		mensagemText := comandosBot(update.Message.Text)

		if mensagemText == "" {
			continue
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, mensagemText)
			msg.ParseMode = "html"
			bot.Send(msg)
		}
	}
}
