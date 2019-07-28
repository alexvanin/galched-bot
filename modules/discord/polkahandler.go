package discord

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/atomic"
)

type polkaHandler struct {
	polka        [][]byte
	voiceChannel string
	lock         *atomic.Bool
}

func (h *polkaHandler) Signature() string {
	return "!song"
}

func (h *polkaHandler) Description() string {
	return "сыграть гимн галчед (только для избранных)"
}

func (h *polkaHandler) IsValid(msg string) bool {
	return msg == "!song"
}

func (h *polkaHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Username != "YoMedved" && m.Author.Username != "Rummy_Quamox" && m.Author.Username != "Lidiya_owl" {
		log.Printf("discord: unathorized polka message from %s", m.Author.Username)
		return
	}

	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return
	}

	// Look for the message sender in that guild's current voice states.
	LogMessage(m)
	if h.lock.CAS(false, true) {
		defer h.lock.Store(false)
		err = h.playSound(s, g.ID, h.voiceChannel)
		if err != nil {
			log.Println("discord: error playing sound:", err)
		}
		time.Sleep(10 * time.Second)
		return
	}
}

// playSound plays the current buffer to the provided channel.
func (h *polkaHandler) playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range h.polka {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specified amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
