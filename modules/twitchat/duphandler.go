package twitchat

import (
	"fmt"
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
)

type (
	dupHandler struct {
		lastMessage string
		counter     int
		dupMinimal  int
	}
)

const DupMinimal = 3

func DupHandler() PrivateMessageHandler {
	return &dupHandler{
		lastMessage: "",
		counter:     0,
		dupMinimal:  DupMinimal,
	}
}

func (h *dupHandler) IsValid(m string) bool {
	return true
}

func (h *dupHandler) Handle(m *twitch.PrivateMessage, r Responser) {
	data := strings.Fields(m.Message)
	for i := range data {
		if data[i] == h.lastMessage {
			h.counter++
		} else {
			if h.counter >= h.dupMinimal {
				msg := fmt.Sprintf("%d %s подряд", h.counter, h.lastMessage)
				r.Say(m.Channel, msg)
				log.Print("chat: ", msg)
			}
			h.counter = 1
			h.lastMessage = data[i]
		}
	}
}
