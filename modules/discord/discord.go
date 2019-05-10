package discord

import (
	"fmt"
	"log"

	"galched-bot/modules/settings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type (
	Discord struct {
		appVersion string
		processor  *HandlerProcessor
		session    *discordgo.Session
	}
)

func New(s *settings.Settings) (*Discord, error) {
	key := fmt.Sprintf("Bot %s", s.DiscordToken)
	instance, err := discordgo.New(key)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create discord instance")
	}

	processor := NewProcessor(s.Version)

	log.Printf("discord: added %d message handlers", len(processor.handlers))
	if len(processor.handlers) > 0 {
		for i := range processor.handlers {
			log.Printf("discord: %d) %s", i+1, processor.handlers[i].Signature())
		}
	}

	return &Discord{
		appVersion: s.Version,
		processor:  processor,
		session:    instance,
	}, nil
}

func LogMessage(m *discordgo.MessageCreate) {
	log.Printf("discord: msg [%s]: %s", m.Author.Username, m.Content)
}

func (d *Discord) Start() error {
	d.session.AddHandler(d.processor.Process)
	return d.session.Open()
}

func (d *Discord) Stop() error {
	return d.session.Close()
}
