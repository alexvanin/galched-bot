package main

import (
	"context"
	"log"

	"go.uber.org/fx"

	"galched-bot/modules/discord"
	"galched-bot/modules/grace"
	"galched-bot/modules/patpet"
	"galched-bot/modules/settings"
	"galched-bot/modules/subday"
	"galched-bot/modules/twitchat"
	"galched-bot/modules/web"
	"galched-bot/modules/youtube"
)

type (
	silentPrinter struct{}

	appParam struct {
		fx.In

		Context  context.Context
		Discord  *discord.Discord
		Settings *settings.Settings
		Chat     *twitchat.TwitchIRC
		Server   *web.WebServer
	}
)

func (s *silentPrinter) Printf(str string, i ...interface{}) {}

func start(p appParam) error {
	var err error

	log.Print("main: starting galched-bot v", p.Settings.Version)

	err = p.Discord.Start()
	if err != nil {
		log.Print("discord: cannot start instance", err)
		return err
	}
	log.Printf("main: discord instance running")

	err = p.Chat.Start()
	if err != nil {
		log.Print("chat: cannot start instance", err)
		return err
	}
	log.Printf("main: twitch chat instance running")

	err = p.Server.Start()
	if err != nil {
		log.Print("web: cannot start instance", err)
		return err
	}
	log.Printf("main: web server instance running")

	log.Printf("main: — — —")
	<-p.Context.Done()
	log.Print("main: stopping galched-bot")

	err = p.Server.Stop(p.Context)
	if err != nil {
		log.Print("web: cannot stop instance", err)
		return err
	}

	err = p.Chat.Stop()
	if err != nil {
		log.Print("chat: cannot stop instance", err)
		return err
	}

	err = p.Discord.Stop()
	if err != nil {
		log.Print("discord: cannot stop instance", err)
		return err
	}

	log.Print("main: galched bot successfully stopped")
	return nil
}

func main() {
	var err error
	app := fx.New(
		fx.Logger(new(silentPrinter)),
		fx.Provide(settings.New, grace.New, discord.New, subday.New,
			twitchat.New, web.New, youtube.New, patpet.New),
		fx.Invoke(start))

	err = app.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
