package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"investohubBot/database"
	"investohubBot/logger"
	"log"
)

type Server struct {
	Bot *tgbotapi.BotAPI
	Db  *database.DBAccess
}

func InitializeBot(bot *tgbotapi.BotAPI, db *database.DBAccess, logger *logger.Logger) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println("Error when listening to updates: ", err)
		logger.MakeLog("Error when listening to updates: " + err.Error())
	}

	for update := range updates {

		if update.CallbackQuery != nil {
			HandleKey(update, bot, db)
		}

		if update.Message != nil {
			HandleCommand(update, bot, db)
		}
	}
}

