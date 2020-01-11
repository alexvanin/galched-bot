package discord

import (
	"fmt"
	"log"
	"strings"

	"galched-bot/modules/subday"

	"github.com/bwmarrin/discordgo"
)

type SubdayListHandler struct {
	subday *subday.Subday
}

func (h *SubdayListHandler) Signature() string {
	return "!sublist"
}
func (h *SubdayListHandler) Description() string {
	return "список игр для сабдея"
}
func (h *SubdayListHandler) IsValid(msg string) bool {
	return msg == "!sublist"
}
func (h *SubdayListHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	h.subday.RLock()
	defer h.subday.RUnlock()
	LogMessage(m)

	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Print("discord: cannot obtain state", err)
		return
	}
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		log.Print("discord: cannot obtain guild", err)
		return
	}
	message := "Игры предыдущих сабдеев доступны по команде **!subhistory**\n" +
		"Список игр для следующего сабдея:\n"
	for k, v := range h.subday.Database() {
		nickname := " "
		for _, member := range g.Members {
			if k == member.User.ID {
				if member.Nick != "" {
					nickname = member.Nick
				} else {
					nickname = member.User.Username
				}
			}
		}
		for _, game := range v {
			message += fmt.Sprintf("   **- %s** от _%s_\n", game, nickname)
		}
	}
	message += "\nВсе команды бота: !galched\n"
	SendMessage(s, m, strings.Trim(message, "\n"))
}

type SubdayAddHandler struct {
	subday *subday.Subday
	roles  []string
}

func (h *SubdayAddHandler) Signature() string {
	return "!subday <game-name>"
}
func (h *SubdayAddHandler) Description() string {
	return "добавление игры в список сабдея"
}
func (h *SubdayAddHandler) IsValid(msg string) bool {
	return strings.HasPrefix(msg, "!subday")
}
func (h *SubdayAddHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	h.subday.Lock()
	defer h.subday.Unlock()
	LogMessage(m)

	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Print("discord: cannot obtain state", err)
		return
	}
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		log.Print("discord: cannot obtain guild", err)
		return
	}
	member, err := s.State.Member(g.ID, m.Author.ID)
	if err != nil {
		log.Print("discord: cannot obtain user role", err)
		return
	}

	if g.Name != "AV" && g.Name != "Galched" {
		log.Printf("discord: message from unsupported guild %s, ignore", g.Name)
		return
	}

	permissionGranted := false
loop:
	for i := range member.Roles {
		for j := range h.roles {
			if member.Roles[i] == h.roles[j] {
				permissionGranted = true
				break loop
			}
		}
	}

	if permissionGranted {
		game := strings.Trim(strings.Replace(m.Content, "!subday", "", 1), " ")
		if game != "" {
			gameList, ok := h.subday.Database()[m.Author.ID]
			if ok && len(gameList) > 10 {
				SendMessage(s, m, "Нельзя заказать больше 10 игр")
				return
			} else if ok {
				for i := range gameList {
					if game == gameList[i] {
						SendMessage(s, m, "Эта игра уже заказана вами")
						return
					}
				}
			}
			h.subday.Database()[m.Author.ID] = append(h.subday.Database()[m.Author.ID], game)
			log.Printf("subday: game [%s] is added to subday database", game)
			SendMessage(s, m, fmt.Sprintf("Игра \"%s\" добавлена в список", game))
			h.subday.DumpToFile()
		}
	} else {
		log.Print("subday: game is not added, insufficient rights")
		SendMessage(s, m, "Заказ игр для сабдея доступен только для подписчиков канала (и Нифлая)")
	}
}

type SubdayHistoryHandler struct{}

func (h *SubdayHistoryHandler) Signature() string {
	return "!subhistory"
}
func (h *SubdayHistoryHandler) Description() string {
	return "история прошлых сабдеев"
}
func (h *SubdayHistoryHandler) IsValid(msg string) bool {
	return strings.HasPrefix(msg, "!subhistory")
}
func (h *SubdayHistoryHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	LogMessage(m)
	message := "Игры предыдущих сабдеев:\n**20.10.18**: _DmC_ -> _Fable 1_ -> _Overcooked 2_\n" +
		"**17.11.18**: _The Witcher_ -> _Xenus: Белое Золото_ -> _NFS: Underground 2_\n" +
		"**22.12.18**: _True Crime: Streets of LA_ -> _Serious Sam 3_ -> _Kholat_\n" +
		"**26.01.19**: _Disney’s Aladdin_ -> _~~Gothic~~_ -> _Scrapland_ -> _Donut County_\n" +
		"**24.02.19**: _Tetris 99_ -> _~~Bully~~_ -> _~~GTA: Vice City~~_\n" +
		"**02.06.19**: _Spec Ops: The Line_ -> _Escape from Tarkov_\n" +
		"**28.07.19**: _Crypt of the Necrodancer_ -> _My Friend Pedro_ -> _Ape Out_\n" +
		"\nВсе команды бота: !galched\n"
	SendMessage(s, m, message)
}

func SubdayHandlers(s *subday.Subday, r []string) []MessageHandler {
	var result []MessageHandler

	addHandler := &SubdayAddHandler{s, r}
	listHandler := &SubdayListHandler{s}
	histHandler := new(SubdayHistoryHandler)
	return append(result, addHandler, listHandler, histHandler)
}
