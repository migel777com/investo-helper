package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func GetMyQuestions(chatID int64, server Server) {
	questionList := server.Db.GetMyQuestions(chatID)

	if len(questionList) == 0 {
		msg := tgbotapi.NewMessage(chatID, "Вы еще не задавали вопросов")
		_, _ = server.Bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Ваши вопросы:")
	_, _ = server.Bot.Send(msg)

	for _, q := range questionList {
		text := "Вопрос: " + q.Question + "\n" +
			"Ответ: "

		if q.Answer == "" {
			text += "На этот вопрос еще не ответили"
		} else {
			text += q.Answer
		}

		msg = tgbotapi.NewMessage(chatID, text)
		_, _ = server.Bot.Send(msg)
	}
}

func GetMasterQuestions(server Server) {
	questionList := server.Db.GetMaserQuestions()

	if len(questionList) == 0 {
		msg := tgbotapi.NewMessage(masterID, "Вопрсов нет")
		_, _ = server.Bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(masterID, "Вопросы:")
	_, _ = server.Bot.Send(msg)

	for _, q := range questionList {
		text := "Вопрос: " + q.Question

		answerCommand := "/question " + strconv.Itoa(int(q.ID))

		msg = tgbotapi.NewMessage(masterID, text)

		answerKey := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Ответить на вопрос", answerCommand),
			),
		)
		msg.ReplyMarkup = answerKey
		_, _ = server.Bot.Send(msg)
	}

}

func AskQuestion(chatID int64, question string, server Server) {
	id := server.Db.AddNewQuestion(chatID, question)
	answerCommand := "/question " + strconv.Itoa(int(id))

	msg := tgbotapi.NewMessage(masterID, "У вас новый вопрос:\n"+question)
	answerKey := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ответить на вопрос", answerCommand),
		),
	)
	msg.ReplyMarkup = answerKey
	_, _ = server.Bot.Send(msg)

	msg = tgbotapi.NewMessage(chatID, "Ваш вопрос принят")
	_, _ = server.Bot.Send(msg)
}
