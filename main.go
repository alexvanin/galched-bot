package main

import (
	"context"
	"log"

	"galched-bot/modules/discord"
	"galched-bot/modules/grace"
	"galched-bot/modules/settings"

	"go.uber.org/fx"
)

type (
	silentPrinter struct{}

	appParam struct {
		fx.In

		Context  context.Context
		Discord  *discord.Discord
		Settings *settings.Settings
	}
)

func (s *silentPrinter) Printf(str string, i ...interface{}) {}

func start(p appParam) {
	var err error

	log.Print("main: starting galched-bot v", p.Settings.Version)

	err = p.Discord.Start()
	if err != nil {
		log.Fatal("discord: cannot start instance", err)
	}
	log.Printf("main: discord instance running")
	log.Printf("main: — — —")

	<-p.Context.Done()
	log.Print("main: stopping galched-bot")

	err = p.Discord.Stop()
	if err != nil {
		log.Fatal("discord: cannot stop instance", err)
	}

	log.Print("main: galched bot successfully stopped")
}

func main() {
	var err error
	app := fx.New(
		fx.Logger(new(silentPrinter)),
		fx.Provide(settings.New, grace.New, discord.New),
		fx.Invoke(start))

	err = app.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
