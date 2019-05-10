package settings

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

const (
	version          = "3.0.0"
	discordTokenPath = "./tokens/.discordtoken"
)

type (
	Settings struct {
		Version      string
		DiscordToken string
	}
)

func New() (*Settings, error) {
	log.Print(os.Getwd())
	discordToken, err := ioutil.ReadFile(discordTokenPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read discord token file")
	}

	return &Settings{
		Version:      version,
		DiscordToken: string(discordToken),
	}, nil
}
