package grifts

import (
	"github.com/StarsiegePlayers/api/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
