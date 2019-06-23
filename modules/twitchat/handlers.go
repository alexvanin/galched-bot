package twitchat

import (
	"github.com/gempir/go-twitch-irc"
)

type (
	MessageHandler interface {
		IsValid(string) bool
		Handle(ch string, u *twitch.User, m *twitch.Message, client *twitch.Client)
	}
)
