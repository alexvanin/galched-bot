package settings

const (
	version          = "3.0.0"
	discordTokenPath = "./tokens/.discordtoken"
)

type (
	Settings struct {
		Version          string
		DiscordTokenPath string
	}
)

func NewSettings() Settings {
	return Settings{
		Version:          version,
		DiscordTokenPath: discordTokenPath,
	}
}
