package twitchat

import (
	"github.com/gempir/go-twitch-irc/v2"
)

type (
	Responser interface {
		Say(channel, message string)
	}

	PrivateMessageHandler interface {
		IsValid(string) bool
		Handle(m *twitch.PrivateMessage, r Responser)
	}
)
