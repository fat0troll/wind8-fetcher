package api

import (
	// local
	"lab.pztrn.name/fat0troll/wind8_fetcher/lib/appcontext"

	// actions
	"lab.pztrn.name/fat0troll/wind8_fetcher/api/fetch_request"
)

var (
	c *appcontext.Context
)

// Initialize prepares API endpoints to initialization
func Initialize(ac *appcontext.Context) {
	c = ac
}

// InitializeEndpoints initializes API endpoints
func InitializeEndpoints() {
	fetchrequest.Initialize(c)
}
