package main

import (
	"context"
	"log"

	"galched-bot/modules/grace"
	"galched-bot/modules/settings"

	"go.uber.org/fx"
)

type (
	silentPrinter struct{}

	appParam struct {
		fx.In

		Context  context.Context
		Settings settings.Settings
	}
)

func (s *silentPrinter) Printf(str string, i ...interface{}) {}

func start(p appParam) {
	log.Print("main: starting galched-bot v", p.Settings.Version)
	<-p.Context.Done()
}

func main() {
	var err error
	app := fx.New(
		fx.Logger(new(silentPrinter)),
		fx.Provide(settings.NewSettings, grace.NewGracefulContext),
		fx.Invoke(start))

	err = app.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
