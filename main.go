package main

import (
	"context"
	"log"

	"galched-bot/modules/discord"
	"galched-bot/modules/grace"
	"galched-bot/modules/settings"
	"galched-bot/modules/subday"

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

func start(p appParam) error {
	var err error

	log.Print("main: starting galched-bot v", p.Settings.Version)

	err = p.Discord.Start()
	if err != nil {
		log.Print("discord: cannot start instance", err)
		return err

	}
	log.Printf("main: discord instance running")
	log.Printf("main: — — —")

	<-p.Context.Done()
	log.Print("main: stopping galched-bot")

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
		fx.Provide(settings.New, grace.New, discord.New, subday.New),
		fx.Invoke(start))

	err = app.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
