package database

import (
	"fmt"
	"investohubBot/models"
	"log"
	"time"
)

func (db *DBAccess) AddNewUser(chatId int64, firstName, lastName, tgUsername, submission string) {
	stmt := "INSERT INTO users(chat_id, first_name, last_name, telegram_username, submission, created_at) VALUES (?, ?, ?, ?, ?, ?);"
	_, err := db.Db.Exec(stmt, chatId, firstName, lastName, tgUsername, submission, time.Now())
	if err != nil {
		return
	}
}

func (db *DBAccess) DeleteUser(chatId int64) {
	stmt := "DELETE FROM users WHERE chat_id = ?"
	db.Db.Exec(stmt, chatId)
}

func (db *DBAccess) CheckUser(chatID int64) bool {
	stmt := "SELECT * FROM users WHERE chat_id = ?;"
	results, err := db.Db.Query(stmt, chatID)

	defer results.Close()

	if err != nil {
		log.Println("Error occurred while checking user existence from database ->", err)
		db.logger.MakeLog("Error occurred while checking user existence from database ->" + err.Error())
	}
	if results.Next() {
		return true
	}
	return false
}

func (db *DBAccess) CheckSubmissionTime(chatID int64) bool {
	stmt := "SELECT created_at FROM users WHERE chat_id = ?;"
	results, err := db.Db.Query(stmt, chatID)

	defer results.Close()

	if err != nil {
		log.Println("Error occurred while checking user submission time from database ->", err)
		db.logger.MakeLog("Error occurred while checking user submission time from database ->" + err.Error())
	}

	var t time.Time

	if results.Next() {
		err := results.Scan(&t)
		if err != nil {
			fmt.Println(err)
		}

		if time.Now().Sub(t).Hours()/24 > 30 {
			return false
		}
	}
	return true
}

func (db *DBAccess) GetMyQuestions(chatID int64)  []models.Question {
	stmt := "SELECT question, answer FROM questions WHERE chat_id = ?;"
	results, err := db.Db.Query(stmt, chatID)

	defer results.Close()

	if err!=nil{
		log.Println("Error occurred while getting single question from database ->", err)
		db.logger.MakeLog("Error occurred while getting single question from database ->"+err.Error())
	}

	var qArray []models.Question

	for results.Next() {
		var question models.Question

		err = results.Scan(&question.Question, &question.Answer)
		if err != nil {
			log.Println("Error occurred while assigning values from db to questions model ->", err)
			db.logger.MakeLog("Error occurred while assigning values from db to questions model ->"+err.Error())
		}

		qArray = append(qArray, question)
	}

	return qArray
}

func (db *DBAccess) GetMaserQuestions()  []models.Question {
	stmt := "SELECT id, question, answer FROM questions WHERE status = 'Not answered';"
	results, err := db.Db.Query(stmt)

	defer results.Close()

	if err!=nil{
		log.Println("Error occurred while getting single question from database ->", err)
		db.logger.MakeLog("Error occurred while getting single question from database ->"+err.Error())
	}

	var qArray []models.Question

	for results.Next() {
		var question models.Question

		err = results.Scan(&question.ID, &question.Question, &question.Answer)
		if err != nil {
			log.Println("Error occurred while assigning values from db to questions model ->", err)
			db.logger.MakeLog("Error occurred while assigning values from db to questions model ->"+err.Error())
		}

		qArray = append(qArray, question)
	}

	return qArray
}

func (db *DBAccess) AddNewQuestion(chatId int64, question string) int64 {
	stmt := "INSERT INTO questions(chat_id, question, answer) VALUES (?, ?, ?);"
	result, err := db.Db.Exec(stmt, chatId, question, "")
	if err != nil {
		log.Println("Error occurred while assigning values from db to questions model ->", err)
		return -1
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error occurred while assigning values from db to questions model ->", err)
		return -1
	}

	return id
}

func (db *DBAccess) CheckInProgress() models.Question {
	stmt := "SELECT id, chat_id, question, answer FROM questions WHERE status='In progress';"
	results, err := db.Db.Query(stmt)

	defer results.Close()

	if err != nil {
		log.Println("Error occurred while checking questions from database ->", err)
		db.logger.MakeLog("Error occurred while checking questions from database ->" + err.Error())
	}

	var question models.Question

	if results.Next() {
		results.Scan(&question.ID, &question.ChatID, &question.Question, &question.Answer)
		question.Status = "In progress"
		return question
	}

	return question
}

func (db *DBAccess) CheckQuestion(id int64) models.Question {
	stmt := "SELECT id, chat_id, question, answer, status FROM questions WHERE id = ?;"
	results, err := db.Db.Query(stmt, id)

	defer results.Close()

	if err != nil {
		log.Println("Error occurred while checking questions from database ->", err)
		db.logger.MakeLog("Error occurred while checking questions from database ->" + err.Error())
	}

	var question models.Question

	if results.Next() {
		results.Scan(&question.ID, &question.ChatID, &question.Question, &question.Answer, &question.Status)
		return question
	}

	return question
}

func (db *DBAccess) StartAnswering(id int64) {
	stmt := "UPDATE questions SET status = 'In progress' WHERE id = ?;"
	db.Db.Exec(stmt, id)
}

func (db *DBAccess) AnswerQuestion(id int64, answer string) {
	stmt := "UPDATE questions SET status = 'Answered', answer = ? WHERE id = ?;"
	db.Db.Exec(stmt, answer, id)
}
