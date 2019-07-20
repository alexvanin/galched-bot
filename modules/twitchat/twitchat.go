package twitchat

import (
	"galched-bot/modules/settings"

	"github.com/gempir/go-twitch-irc/v2"
)

type (
	TwitchIRC struct {
		username string
		chat     *twitch.Client
		handlers []PrivateMessageHandler
	}
)

func New(s *settings.Settings) (*TwitchIRC, error) {
	var irc = new(TwitchIRC)

	irc.username = s.TwitchUser

	irc.handlers = append(irc.handlers, DupHandler())

	irc.chat = twitch.NewClient(s.TwitchUser, s.TwitchToken)
	irc.chat.OnPrivateMessage(irc.PrivateMessageHandler)
	irc.chat.Join(s.TwitchIRCRoom)

	return irc, nil
}

func (c *TwitchIRC) Start() error {
	go func() {
		err := c.chat.Connect()
		_ = err // no point in error because disconnect will be called anyway
	}()
	return nil
}

func (c *TwitchIRC) Stop() error {
	return c.chat.Disconnect()
}

func (c *TwitchIRC) PrivateMessageHandler(msg twitch.PrivateMessage) {
	if msg.User.Name == c.username {
		return
	}
	for i := range c.handlers {
		if c.handlers[i].IsValid(msg.Message) {
			c.handlers[i].Handle(&msg, c.chat)
		}
	}
}
