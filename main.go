package main

import (
	"log"
	"os"

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
	utils.Convert()

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

		switch update.Message.Command() {
		case "start":
			message.ReplyMarkup = kb
			message.Text = "Hi !"

		case "help":
			var ansList []interface{}
			files := utils.Convert()
			for _, file := range len(files) {
				ansList = append(ansList, tg_api.NewInputMediaPhoto(tg_api.FilePath((file))))

			}
			//text, _ := utils.GetPdf()
			media_group := tg_api.NewMediaGroup(update.Message.Chat.ID, ansList)
			message.Text = "sd"
			bot.SendMediaGroup(media_group)
			if _, err := bot.Send(message); err != nil {
				log.Fatal(err)
			}
		}
	}

	// fmt.Println(utils.GetPdf())

}
