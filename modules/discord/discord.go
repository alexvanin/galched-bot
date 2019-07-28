package discord

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"galched-bot/modules/settings"
	"galched-bot/modules/subday"

	"go.uber.org/atomic"

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

	polka, err := loadSong(s.PolkaPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read polka song")
	}
	processor.AddHandler(&polkaHandler{
		polka:        polka,
		voiceChannel: s.DiscordVoiceChannel,
		lock:         atomic.NewBool(false),
	})

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

func loadSong(path string) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "error opening dca file")
	}

	var (
		opuslen int16
		buffer  = make([][]byte, 0)
	)

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return nil, err
			}
			return buffer, nil
		}

		if err != nil {
			return nil, errors.Wrap(err, "error reading from dca file")
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			return nil, errors.Wrap(err, "error  reading from dca file")
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
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
