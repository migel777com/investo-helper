package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"investohubBot/database"
	"investohubBot/logger"
	"investohubBot/tgbot"
	"io/ioutil"
	"log"
	"os"
)

var (
	bot *tgbotapi.BotAPI
)

func initTelegram() {
	_ = godotenv.Load("global.env")
	botToken := os.Getenv("BOT_TOKEN")
	var err error

	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Println(err)
		return
	}

	// this perhaps should be conditional on GetWebhookInfo()
	// only set webhook if it is not set properly
	url := "https://investo-helper.herokuapp.com/" + bot.Token
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(url))
	if err != nil {
		log.Println(err)
	}
}

func webhookHandler(c *gin.Context) {
	myLogger := logger.RegisterLogger("logs.txt")
	defer myLogger.File.Close()

	db, _ := database.OpenDB(&myLogger)
	defer db.Db.Close()

	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}

	// to monitor changes run: heroku logs --tail

	if update.CallbackQuery != nil {
		tgbot.HandleKey(update, bot, &db)
		log.Printf("From: %+v Text: %+v ChatID: %+v\n", update.CallbackQuery.From, update.CallbackQuery.Data, update.CallbackQuery.Message.Chat.ID)
	}

	if update.Message != nil {
		tgbot.HandleCommand(update, bot, &db)
		log.Printf("From: %+v Text: %+v ChatID: %+v\n", update.Message.From, update.Message.Text, update.Message.Chat.ID)
	}
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// gin router
	router := gin.New()
	router.Use(gin.Logger())

	// telegram
	initTelegram()
	fmt.Println(bot.Token)
	router.POST("/"+bot.Token, webhookHandler)

	err := router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}

	/*_ = godotenv.Load("global.env")
	botToken := os.Getenv("BOT_TOKEN")

	myLogger := logger.RegisterLogger("logs.txt")
	defer myLogger.File.Close()

	db, _ := database.OpenDB(&myLogger)
	defer db.Db.Close()

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil{
		log.Println("Error when creating a bot instance occurred", err)
	}

	tgbot.InitializeBot(bot, &db, &myLogger)*/
}
