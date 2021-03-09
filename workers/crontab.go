package workers

import (
	"github.com/gobuffalo/buffalo"
	"github.com/robfig/cron/v3"
)

// https://pkg.go.dev/github.com/robfig/cron

var app *buffalo.App

func InitCron(cron *cron.Cron, appIn *buffalo.App) {
	app = appIn
	_, _ = cron.AddFunc("@every 5m", UpdateServerList)

	go UpdateServerList()
}