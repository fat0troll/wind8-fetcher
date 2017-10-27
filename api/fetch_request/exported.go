package fetchrequest

import (
	// local
	"lab.pztrn.name/fat0troll/wind8_fetcher/lib/appcontext"
)

var (
	c *appcontext.Context
)

// Initialize is an initialization function for call request handler
func Initialize(ac *appcontext.Context) {
	c = ac
	c.Log.Info("Initializing action for /fetch route...")

	c.HTTPServerMux.HandleFunc("/fetch", HandleFetchRequest)
}
