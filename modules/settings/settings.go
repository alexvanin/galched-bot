package settings

import (
	"io/ioutil"
	"log"
)

const (
	version          = "3.0.2"
	discordTokenPath = "./tokens/.discordtoken"
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
		SubdayDataPath string
		PermittedRoles []string
	}
)

func New() (*Settings, error) {
	discordToken, err := ioutil.ReadFile(discordTokenPath)
	if err != nil {
		log.Print("settings: cannot read discord token file", err)
	}

	return &Settings{
		Version:        version,
		DiscordToken:   string(discordToken),
		SubdayDataPath: subdayDataPath,
		PermittedRoles: []string{subRole1, subRole2, galchedRole, smorcRole},
	}, nil
}
