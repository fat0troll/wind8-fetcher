package main

import (
	// local
	"source.wtfteam.pro/wind8/wind8_fetcher/api"
	"source.wtfteam.pro/wind8/wind8_fetcher/lib/appcontext"
)

func main() {
	c := appcontext.New()
	c.Init()
	c.InitializeStartupFlags()
	c.StartupFlags.Parse()

	configPath, err := c.StartupFlags.GetStringValue("config")
	if err != nil {
		c.Log.Errorln(err)
		c.Log.Fatal("Can't get config file parameter from command line. Exiting.")
	}
	c.InitializeConfig(configPath)

	c.Log.Info("Starting API endpoints...")
	api.Initialize(c)
	api.InitializeEndpoints()

	c.StartHTTPListener()
}
