package discord

import (
	"fmt"

	"galched-bot/modules/settings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type (
	Discord struct {
		session *discordgo.Session
	}
)

func New(s *settings.Settings) (*Discord, error) {
	key := fmt.Sprintf("Bot %s", s.DiscordToken)
	instance, err := discordgo.New(key)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create discord instance")
	}
	return &Discord{session: instance}, nil
}

func (d *Discord) Start() error {
	return d.session.Open()
}

func (d *Discord) Stop() error {
	return d.session.Close()
}
