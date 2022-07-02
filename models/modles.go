package models

import (
	"time"
)

type Question struct {
	ID       int64
	ChatID   int64
	Question string
	Answer   string
	Status   string
}

type User struct {
	ID               int64
	ChatID           int64
	FirstName        string
	LastName         string
	TelegramUsername string
	Submission       string
	CreatedAt        time.Time
}

type OperationStateResponse struct {
	Result Result `xml:"Result"`
	State  State  `xml:"State"`
}

type Result struct {
	Code        int64  `xml:"Code"`
	Description string `xml:"Description"`
}

type State struct {
	Code int64 `xml:"Code"`
}
