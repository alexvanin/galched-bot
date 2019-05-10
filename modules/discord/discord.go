package discord

import (
	"fmt"
	"log"

	"galched-bot/modules/settings"
	"galched-bot/modules/subday"

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

func New(s *settings.Settings, subday *subday.Subday) (*Discord, error) {
	key := fmt.Sprintf("Bot %s", s.DiscordToken)
	instance, err := discordgo.New(key)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create discord instance")
	}

	processor := NewProcessor(s.Version)
	for _, subdayHandler := range SubdayHandlers(subday, s.PermittedRoles) {
		processor.AddHandler(subdayHandler)
	}

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

func SendMessage(s *discordgo.Session, m *discordgo.MessageCreate, text string) {
	_, err := s.ChannelMessageSend(m.ChannelID, text)
	if err != nil {
		log.Printf("discord: cannot send message [%s]: %v", text, err)
	}
}

func (d *Discord) Start() error {
	d.session.AddHandler(d.processor.Process)
	return d.session.Open()
}

func (d *Discord) Stop() error {
	return d.session.Close()
}
