package twitchat

import (
	"fmt"
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc"
)

type (
	dupHandler struct {
		lastMessage string
		counter     int
		dupMinimal  int
	}
)

const DupMinimal = 3

func DupHandler() MessageHandler {
	return &dupHandler{
		lastMessage: "",
		counter:     0,
		dupMinimal:  DupMinimal,
	}
}

func (h *dupHandler) IsValid(m string) bool {
	return true
}

func (h *dupHandler) Handle(ch string, u *twitch.User, m *twitch.Message, client *twitch.Client) {
	data := strings.Fields(m.Text)
	for i := range data {
		if data[i] == h.lastMessage {
			h.counter++
		} else {
			if h.counter >= h.dupMinimal {
				msg := fmt.Sprintf("%d %s подряд", h.counter, h.lastMessage)
				client.Say(ch, msg)
				log.Print("chat: ", msg)
			}
			h.counter = 1
			h.lastMessage = data[i]
		}
	}
}
