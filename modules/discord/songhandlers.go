package discord

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"go.uber.org/atomic"

	"galched-bot/modules/settings"
)

type (
	songData [][]byte

	SongHandler struct {
		globalLock *atomic.Bool
		songLock   *atomic.Bool

		song         songData
		signature    string
		description  string
		voiceChannel string
		permissions  []string
		timeout      time.Duration
	}
)

func (h *SongHandler) Signature() string {
	return h.signature
}

func (h *SongHandler) Description() string {
	return h.description
}

func (h *SongHandler) IsValid(msg string) bool {
	return msg == h.signature
}

func (h *SongHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	var permitted bool
	for i := range h.permissions {
		if m.Author.Username == h.permissions[i] {
			permitted = true
			break
		}
	}
	if len(h.permissions) > 0 && !permitted {
		log.Printf("discord: unathorized %s message from %s",
			h.signature, m.Author.Username)
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
	if h.globalLock.CAS(false, true) {
		if h.songLock.CAS(false, true) {

			err = playSound(s, g.ID, h.voiceChannel, h.song)
			if err != nil {
				log.Println("discord: error playing sound:", err)
			}

			h.globalLock.Store(false)
			time.Sleep(h.timeout)
			defer h.songLock.Store(false)

			return
		}
		h.globalLock.Store(false)
	}
}

func SongHandlers(s *settings.Settings) []MessageHandler {
	result := make([]MessageHandler, 0, len(s.Songs))
	g := new(atomic.Bool)

	for i := range s.Songs {
		song, err := loadSong(s.Songs[i].Path)
		if err != nil {
			log.Println("discord: error loading song file", err)
			continue
		}
		handler := &SongHandler{
			globalLock:   g,
			songLock:     new(atomic.Bool),
			song:         song,
			signature:    s.Songs[i].Signature,
			description:  s.Songs[i].Description,
			voiceChannel: s.DiscordVoiceChannel,
			permissions:  s.Songs[i].Permissions,
			timeout:      s.Songs[i].Timeout,
		}
		result = append(result, handler)
	}

	return result
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string, song songData) error {

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
	for _, buff := range song {
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

func loadSong(path string) (songData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "error opening dca file")
	}

	var (
		opuslen int16
		buffer  = make(songData, 0)
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
