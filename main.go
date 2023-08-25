package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"utils/utils"

	"github.com/joho/godotenv"

	tg_api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var kb = tg_api.NewReplyKeyboard(
	tg_api.NewKeyboardButtonRow(
		tg_api.NewKeyboardButton("Замены"),
		tg_api.NewKeyboardButton("Цвет недели"),
	),
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env not found")
	}
	token, _ := os.LookupEnv("API_TOKEN")

	bot, err := tg_api.NewBotAPI(token)
	if err != nil {
		log.Fatal(err, " Token not found")
	}

	bot.Debug = true

	log.Printf("Authorization on acc %s", bot.Self.UserName)

	u := tg_api.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
		fmt.Println(string(update.Message.Text))
		start := time.Now()
		switch string(update.Message.Text) {
		case "/start":
			message := tg_api.NewMessage(update.Message.Chat.ID, "Hi")
			message.ReplyMarkup = kb
			bot.Send(message)

		case "/schedule":
			start := time.Now()
			result := []interface{}{}
			fileChan := make(chan tg_api.InputMediaPhoto)
			c2 := make(chan string)
			go utils.GetPdf(c2)
			go utils.Convert(fileChan, <-c2)

			//text, _ := utils.GetPdf()
			for {
				if val, opened := <-fileChan; opened {
					result = append(result, val)
				} else {
					break
				}
			}
			media_group := tg_api.NewMediaGroup(update.Message.Chat.ID, result)
			fmt.Println(time.Since(start))
			bot.Send(media_group)

		}
		fmt.Print(time.Since(start))
	}

}
