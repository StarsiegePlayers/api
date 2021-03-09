package actions

import (
	"github.com/StarsiegePlayers/api/workers"
	"github.com/gobuffalo/buffalo"
	"net/http"
)

// ServerListHandler is a handler to serve up
// the multiplayer server list
func ServerListHandler(c buffalo.Context) error {

	return c.Render(http.StatusOK, r.JSON(workers.GetServerList()))
}

