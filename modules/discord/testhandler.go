package discord

import (
	"github.com/bwmarrin/discordgo"
)

type testHandler struct{}

var _ = (testHandler)(nil) // ignore unused warning

func (t *testHandler) Signature() string {
	return "!test"
}

func (t *testHandler) Description() string {
	return "тестовый хэндлер"
}

func (t *testHandler) IsValid(msg string) bool {
	return msg == "!test"
}

func (t *testHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	LogMessage(m)
}
