package config

import (
	"github.com/george012/gtbox"
	"github.com/george012/gtbox/gtbox_app"
)

var (
	CurrentApp *ExtendApp
)

type ExtendApp struct {
	*gtbox_app.App
	APIPort  int
	MCPPort  int
	GRPCPort int
}

func NewApp(appName, bundleID, description string, runMode gtbox.RunMode) *ExtendApp {
	app := &ExtendApp{
		App: gtbox_app.NewApp(appName, ProjectVersion, bundleID, description, runMode),
	}

	return app
}
