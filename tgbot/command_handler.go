package tgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"investohubBot/database"
	"investohubBot/robokassa"
	"strconv"
	"strings"
)

//var masterID int64 = 622108583 my
var masterID int64 = 1100797127

func HandleCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *database.DBAccess) {
	server := Server{Bot: bot, Db: db}

	chatID := update.Message.Chat.ID
	message := update.Message.Text

	if chatID == masterID {

		switch message {
		case "Вопросы":
			GetMasterQuestions(server)
		default:
			q := server.Db.CheckInProgress()

			if q.Status == "In progress" {
				server.Db.AnswerQuestion(q.ID, message)

				msg := tgbotapi.NewMessage(q.ChatID, "Ответили на ваш вопрос\n"+
					q.Question+"\n "+
					"Ответ на вопрос: "+message)
				msg.ReplyMarkup = mainMenu
				_, _ = server.Bot.Send(msg)
				return
			}

			msg := tgbotapi.NewMessage(masterID, "Welcome Master")
			msg.ReplyMarkup = masterMenu
			_, _ = server.Bot.Send(msg)
		}

		return
	}

	if !server.Db.CheckUser(chatID) {
		if message == "/start" {
			msg := tgbotapi.NewMessage(chatID, "Привет👋🏻\n\n"+
				"Команда Investo Drive создала для тебя консультационный бот, в котором ты сможешь задать любой вопрос по "+
				"инвестициям нашим менторам. \n\n"+
				"🔐 Оплатив подписку в размере 3500 рублей (20 127 тенге) ты получишь "+
				"постоянный доступ к личному кабинету, где ты сможешь задавать анонимные вопросы по инвестициям нашим "+
				"менторам и получать ответ в кратчайшие сроки.\n\n"+
				"✅ Обучайтесь у профессионалов Investo Drive и проверяйте гипотезы вместе с нами!\n\n"+
				"🔗Ссылка на договор оферты: "+
				"https://drive.google.com/file/d/1DTAplzE3RfPSc9pKEEOtZIXm8WkGbMxP/view?usp=sharing")
			msg.ReplyMarkup = payment
			_, _ = server.Bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Подписка будет стоит 50$/мес\n\n"+
				"Пройти по ссылке нажав на кнопку \"Оплатит\". \n\n"+
				"После оплаты нажмите \"Провеерить оплату\", чтобы получить доступ.")
			var pay = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("Оплатить", robokassa.GetURL(chatID)),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Проверить оплату", "/checkPayment"),
				),
			)
			msg.ReplyMarkup = pay
			_, _ = server.Bot.Send(msg)
		}
		fmt.Println(chatID)
		fmt.Println(message)
		return
	}

	if !server.Db.CheckSubmissionTime(chatID) {
		msg := tgbotapi.NewMessage(chatID, "Ваша подписка закончилась")
		msg.ReplyMarkup = payment
		_, _ = server.Bot.Send(msg)

		server.Db.DeleteUser(chatID)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "")
	msg.ReplyMarkup = mainMenu
	_, _ = server.Bot.Send(msg)

	var parameter string

	if strings.Contains(message, ":") {
		array := strings.Split(message, ":")
		if array[0] == "Вопрос" {
			message = array[0]
			for i := 1; i < len(array); i++ {
				parameter += array[i]
			}
		}
	}

	switch message {
	case "Мои вопросы":
		GetMyQuestions(chatID, server)
	case "Задать вопрос":
		msg = tgbotapi.NewMessage(chatID, "Что бы задать вопрос вам необходимо написать кодовое слово и символ двоеточия "+
			"(Рекомендуется скопировать и вставить перед текстом вопроса!) \"Вопрос:\". "+
			"ВНИМАНИЕ: Вопросы не составленные по вышеописанному формату рассматриваться не будут.")
		_, _ = server.Bot.Send(msg)
	case "Вопрос":
		AskQuestion(chatID, parameter, server)
	default:
		msg = tgbotapi.NewMessage(chatID, "Выберете действие из меню")
		msg.ReplyMarkup = mainMenu
		_, _ = server.Bot.Send(msg)
	}

}

func HandleKey(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *database.DBAccess) {
	server := Server{Bot: bot, Db: db}

	chatID := update.CallbackQuery.Message.Chat.ID
	message := update.CallbackQuery.Data

	var parameter string

	if strings.Contains(message, " ") {
		array := strings.Split(message, " ")
		if array[0] == "/question" {
			message = array[0]
			parameter = array[1]
		}
	}

	if chatID == masterID {
		switch message {
		case "/question":
			q := server.Db.CheckInProgress()

			fmt.Println(q.Status)

			if q.Status == "In progress" {
				msg := tgbotapi.NewMessage(chatID, "Вы еще не ответили на этот вопрос:\n"+q.Question)
				_, _ = server.Bot.Send(msg)
				return
			}

			id, err := strconv.ParseInt(parameter, 10, 64)
			if err != nil {
				fmt.Println(err)
				return
			}

			q = server.Db.CheckQuestion(id)
			if q.Status == "Answered" {
				msg := tgbotapi.NewMessage(chatID, "Вы уже ответили на этот вопрос")
				_, _ = server.Bot.Send(msg)
				return
			}

			msg := tgbotapi.NewMessage(chatID, "Дайте ваш ответ на вопрос "+q.Question)
			_, _ = server.Bot.Send(msg)

			server.Db.StartAnswering(id)
		}

		return
	}

	if !server.Db.CheckUser(chatID) {
		if message == "/payment" {
			msg := tgbotapi.NewMessage(chatID, "Подписка будет стоить 3500 рублей (20 127 тенге).\n\n"+
				"Оплатить ты сможешь нажав на кнопку \"Оплатить\" внизу.\n\n"+
				"После оплаты, нажми на \"Проверить оплату\", чтобы открыть доступ к Телеграм боту.")
			var pay = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("Оплатить", robokassa.GetURL(chatID)),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Проверить оплату", "/checkPayment"),
				),
			)
			msg.ReplyMarkup = pay
			_, _ = server.Bot.Send(msg)
		} else if message == "/checkPayment" {
			//fmt.Println("Проверить оплату")

			status := robokassa.CheckPayment(chatID)

			if status == 100 {
				msg := tgbotapi.NewMessage(chatID, "Ваша оплата успешно получена!")
				msg.ReplyMarkup = mainMenu
				server.Db.AddNewUser(chatID, update.CallbackQuery.From.FirstName, update.CallbackQuery.From.LastName, update.CallbackQuery.From.UserName, "premium")
				_, _ = server.Bot.Send(msg)
			} else if status == 50 {
				msg := tgbotapi.NewMessage(chatID, "Ваш платеж находится в обработке...")
				_, _ = server.Bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(chatID, "Мы еще не получили вашу оплату.")
				_, _ = server.Bot.Send(msg)
			}

			//server.Db.AddNewUser(chatID, update.CallbackQuery.From.FirstName, update.CallbackQuery.From.LastName, update.CallbackQuery.From.UserName, "premium")
		}

		/*fmt.Println(chatID)
		fmt.Println(message)*/
		return
	}

	if !server.Db.CheckSubmissionTime(chatID) {
		msg := tgbotapi.NewMessage(chatID, "Ваша подписка закончилась")
		msg.ReplyMarkup = payment
		_, _ = server.Bot.Send(msg)

		server.Db.DeleteUser(chatID)
		return
	}

}
