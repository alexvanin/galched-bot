package settings

import (
	"io/ioutil"
	"log"
)

const (
	version          = "4.0.0"
	twitchUser       = "galchedbot"
	twitchIRCRoom    = "galched"
	discordTokenPath = "./tokens/.discordtoken"
	twitchTokenPath  = "./tokens/.twitchtoken"
	subdayDataPath   = "./backups/subday"

	// Permitted roles in discord for subday
	subRole1    = "433672344737677322"
	subRole2    = "433680494635515904"
	galchedRole = "301467455497175041"
	smorcRole   = "301470784491356172"
)

type (
	Settings struct {
		Version        string
		DiscordToken   string
		TwitchUser     string
		TwitchIRCRoom  string
		TwitchToken    string
		SubdayDataPath string
		PermittedRoles []string
	}
)

func New() (*Settings, error) {
	discordToken, err := ioutil.ReadFile(discordTokenPath)
	if err != nil {
		log.Print("settings: cannot read discord token file", err)
	}
	twitchToken, err := ioutil.ReadFile(twitchTokenPath)
	if err != nil {
		log.Print("settings: cannot read twitch token file", err)
	}

	return &Settings{
		Version:        version,
		DiscordToken:   string(discordToken),
		TwitchToken:    string(twitchToken),
		TwitchUser:     twitchUser,
		TwitchIRCRoom:  twitchIRCRoom,
		SubdayDataPath: subdayDataPath,
		PermittedRoles: []string{subRole1, subRole2, galchedRole, smorcRole},
	}, nil
}
