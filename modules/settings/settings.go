package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

const (
	version          = "5.0.1"
	twitchUser       = "galchedbot"
	twitchIRCRoom    = "galched"
	discordTokenPath = "./tokens/.discordtoken"
	twitchTokenPath  = "./tokens/.twitchtoken"
	subdayDataPath   = "./backups/subday"
	youtubeTokenPath = "./tokens/.youtubetoken"
	webLoginsPath    = "./tokens/.weblogins"

	// Permitted roles in discord for subday
	subRole1    = "433672344737677322"
	subRole2    = "433680494635515904"
	galchedRole = "301467455497175041"
	smorcRole   = "301470784491356172"

	defaultQueueAddr = ":8888"
)

type (
	SongInfo struct {
		Path        string
		Signature   string
		Description string
		Permissions []string
		Timeout     time.Duration
	}

	Settings struct {
		Version             string
		DiscordToken        string
		TwitchUser          string
		TwitchIRCRoom       string
		TwitchToken         string
		YoutubeToken        string
		SubdayDataPath      string
		PermittedRoles      []string
		DiscordVoiceChannel string
		Songs               []SongInfo

		QueueAddress string
		LoginUsers   map[string]string
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
	youtubetoken, err := ioutil.ReadFile(youtubeTokenPath)
	if err != nil {
		log.Print("settings: cannot read twitch token file", err)
	}

	webLogins := make(map[string]string)
	webLoginsRaw, err := ioutil.ReadFile(webLoginsPath)
	if err != nil {
		log.Print("settings: cannot read web login file", err)
	} else {
		err = json.Unmarshal(webLoginsRaw, &webLogins)
		if err != nil {
			log.Print("settings: cannot parse web login file", err)
		}
	}

	return &Settings{
		Version:             version,
		DiscordToken:        string(discordToken),
		TwitchToken:         string(twitchToken),
		YoutubeToken:        string(youtubetoken),
		TwitchUser:          twitchUser,
		TwitchIRCRoom:       twitchIRCRoom,
		SubdayDataPath:      subdayDataPath,
		DiscordVoiceChannel: "301793085522706432",
		PermittedRoles:      []string{subRole1, subRole2, galchedRole, smorcRole},
		Songs: []SongInfo{
			{
				Path:        "songs/polka.dca",
				Signature:   "!song",
				Description: "сыграть гимн галчед (только для избранных)",
				Permissions: []string{"AlexV", "Rummy_Quamox", "Lidiya_owl"},
				Timeout:     10 * time.Second,
			},
			{
				Path:        "songs/whisper.dca",
				Signature:   "!sax",
				Description: "kreygasm",
				Timeout:     20 * time.Second,
			},
			{
				Path:        "songs/st.dca",
				Signature:   "!chiki",
				Description: "briki v damki",
				Timeout:     20 * time.Second,
			},
		},
		QueueAddress: defaultQueueAddr,
		LoginUsers:   webLogins,
	}, nil
}
