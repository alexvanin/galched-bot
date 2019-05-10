package discord

import (
	"github.com/bwmarrin/discordgo"
)

type testHandler struct{}

var _ = testHandler{} // ignore unused warning

func (h *testHandler) Signature() string {
	return "!test"
}

func (h *testHandler) Description() string {
	return "тестовый хэндлер"
}

func (h *testHandler) IsValid(msg string) bool {
	return msg == "!test"
}

func (h *testHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	LogMessage(m)
}
