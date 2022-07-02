package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var payment = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Оплатить","/payment"),
	),
)

var pay = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("Оплатить",""),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Проверить","/checkPayment"),
	),
)

var mainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Мои вопросы"),
		tgbotapi.NewKeyboardButton("Задать вопрос"),
	),
)

var masterMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Вопросы"),
	),
)
