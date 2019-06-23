package twitchat

import (
	"galched-bot/modules/settings"

	"github.com/gempir/go-twitch-irc"
	_ "github.com/gempir/go-twitch-irc"
)

type (
	TwitchIRC struct {
		username string
		chat     *twitch.Client
		handlers []MessageHandler
	}
)

func New(s *settings.Settings) (*TwitchIRC, error) {
	var irc = new(TwitchIRC)

	irc.username = s.TwitchUser

	irc.handlers = append(irc.handlers, DupHandler())

	irc.chat = twitch.NewClient(s.TwitchUser, s.TwitchToken)
	irc.chat.OnNewMessage(irc.MessageHandler)
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

func (c *TwitchIRC) MessageHandler(ch string, u twitch.User, m twitch.Message) {
	if u.Username == c.username {
		return
	}
	for i := range c.handlers {
		if c.handlers[i].IsValid(m.Text) {
			c.handlers[i].Handle(ch, &u, &m, c.chat)
		}
	}
}
