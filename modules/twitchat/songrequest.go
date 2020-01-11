package twitchat

import (
	"fmt"
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"

	"galched-bot/modules/youtube"
)

const (
	songMsg   = "!song"
	reqPrefix = "!req " // space in the end is important
)

type (
	songRequest struct {
		r *youtube.Requester
	}
)

func SongRequest(r *youtube.Requester) PrivateMessageHandler {
	return &songRequest{r: r}
}

func (h *songRequest) IsValid(m *twitch.PrivateMessage) bool {
	return (strings.HasPrefix(m.Message, reqPrefix) && m.Tags["msg-id"] == "highlighted-message") ||
		strings.TrimSpace(m.Message) == songMsg
	// return strings.HasPrefix(m.Message, reqPrefix) || strings.TrimSpace(m.Message) == songMsg
}

func (h *songRequest) Handle(m *twitch.PrivateMessage, r Responser) {
	if strings.TrimSpace(m.Message) == "!song" {
		list := h.r.List()
		if len(list) > 0 {
			line := fmt.Sprintf("Сейчас играет: <%s>", list[0].Title)
			r.Say(m.Channel, line)
		} else {
			r.Say(m.Channel, "Очередь видео пуста")
		}
		return
	}

	query := strings.TrimPrefix(m.Message, "!req ")
	if len(query) == 0 {
		return
	}

	chatMsg, err := h.r.AddVideo(query, m.User.DisplayName)
	if err != nil {
		log.Printf("yt: cannot add song from msg <%s>, err: %v", m.Message, err)
		if len(chatMsg) > 0 {
			r.Say(m.Channel, m.User.DisplayName+" "+chatMsg)
		}
		return
	}
	r.Say(m.Channel, m.User.DisplayName+" добавил "+chatMsg)
	return
}
