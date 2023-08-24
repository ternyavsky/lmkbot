package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"utils/utils"

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

		message := tg_api.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Text {
		case "/start":
			start := time.Now()
			message.ReplyMarkup = kb
			message.Text = "Hi !"
			bot.Send(message)
			fmt.Println(time.Since(start))

		case "/shedule":
			start := time.Now()
			var ansList []interface{}
			caption, _ := utils.GetPdf()
			files := utils.Convert()
			for i, file := range files {
				added := tg_api.NewInputMediaPhoto(tg_api.FilePath(file))
				fmt.Println(i, file)
				if i == 1 {
					added.Caption = caption
					ansList = append(ansList, added)
				} else {
					ansList = append(ansList, added)
				}

			}
			//text, _ := utils.GetPdf()
			media_group := tg_api.NewMediaGroup(update.Message.Chat.ID, ansList)
			message.Text = "sd"
			bot.Send(media_group)
			fmt.Print(time.Since(start))
		}
	}

}
