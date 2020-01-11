package twitchat

import (
	"log"

	"github.com/gempir/go-twitch-irc/v2"
)

type (
	logCheck struct{}
)

func LogCheck() PrivateMessageHandler {
	return new(logCheck)
}

func (h *logCheck) IsValid(m *twitch.PrivateMessage) bool {
	return true
}

func (h *logCheck) Handle(m *twitch.PrivateMessage, r Responser) {
	log.Print("chat <", m.User.DisplayName, "> : ", m.Message)
}
