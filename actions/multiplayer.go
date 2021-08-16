package actions

import (
	"net/http"
	"time"

	"github.com/StarsiegePlayers/api/workers"
	"github.com/gobuffalo/buffalo"
)

// ServerListHandler is a handler to serve up
// the multiplayer server list
func ServerListHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON(workers.GetServerList()))
}

func ServerListLastUpdated(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		serverList := workers.GetServerList()
		c.Response().Header().Set("Last-Modified", serverList.RequestTime.Format(time.RFC1123))
		return next(c)
	}
}

func ServerStatistics(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON(workers.GetServerList()))
}

func ServerStatisticsLastUpdated(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		serverList := workers.GetServerList()
		c.Response().Header().Set("Last-Modified", serverList.RequestTime.Format(time.RFC1123))
		return next(c)
	}
}
