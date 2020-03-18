package twitchat

import (
	"fmt"
	"strings"

	"galched-bot/modules/patpet"
	"github.com/gempir/go-twitch-irc/v2"
)

const (
	petMsg1 = "!погладь"
	petMsg2 = "!гладь"
	petMsg3 = "!погладить"
)

type (
	petCat struct {
		cat *patpet.Pet
	}
)

func PetCat(pet *patpet.Pet) *petCat {
	return &petCat{
		cat: pet,
	}
}

func (h *petCat) IsValid(m *twitch.PrivateMessage) bool {
	return (m.Tags["msg-id"] == "highlighted-message") && (strings.HasPrefix(m.Message, petMsg1) ||
		strings.HasPrefix(m.Message, petMsg2) ||
		strings.HasPrefix(m.Message, petMsg3))
}

func (h *petCat) Handle(m *twitch.PrivateMessage, r Responser) {
	msg := fmt.Sprintf("Котэ поглажен уже %d раз(а) InuyoFace", h.cat.Pet())
	r.Say(m.Channel, msg)
	h.cat.Dump()
}
