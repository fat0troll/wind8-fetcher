package appcontext

import (
	// stdlib
	"encoding/json"
	"net/http"
	"os"
	// 3rd-party
	"lab.pztrn.name/golibs/flagger"
	"lab.pztrn.name/golibs/mogrus"
	// local
	"lab.pztrn.name/fat0troll/wind8_fetcher/lib/config"
)

// Context is an application context struct
type Context struct {
	Cfg           *config.Config
	HTTPServerMux *http.ServeMux
	Log           *mogrus.LoggerHandler
	StartupFlags  *flagger.Flagger
}

// Init is an initialization function for context
func (c *Context) Init() {
	l := mogrus.New()
	l.Initialize()
	c.Log = l.CreateLogger("stdout")
	c.Log.CreateOutput("stdout", os.Stdout, true)

	c.Cfg = config.New()

	c.StartupFlags = flagger.New(c.Log)
	c.StartupFlags.Initialize()
	c.HTTPServerMux = http.NewServeMux()
}

// InitializeConfig fills config struct with data from given file
func (c *Context) InitializeConfig(configPath string) {
	c.Cfg.Init(c.Log, configPath)
}

// InitializeStartupFlags gives information about available startup flags
func (c *Context) InitializeStartupFlags() {
	configFlag := flagger.Flag{}
	configFlag.Name = "config"
	configFlag.Description = "Configuration file path"
	configFlag.Type = "string"
	configFlag.DefaultValue = "config.yaml"
	err := c.StartupFlags.AddFlag(&configFlag)
	if err != nil {
		c.Log.Errorln(err)
	}
}

// StartHTTPListener starts HTTP server on given port
func (c *Context) StartHTTPListener() {
	response := make(map[string]string)
	responseBody := make([]byte, 0)
	c.HTTPServerMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(404)
		response["status"] = "error"
		response["error_coder"] = "404"
		response["descirption"] = "Not found."

		responseBody, _ = json.Marshal(response)
		w.Write(responseBody)
	})

	c.Log.Info("HTTP server started at http://" + c.Cfg.HTTPListener.Host + ":" + c.Cfg.HTTPListener.Port)
	err := http.ListenAndServe(c.Cfg.HTTPListener.Host+":"+c.Cfg.HTTPListener.Port, c.HTTPServerMux)
	c.Log.Fatalln(err)
}
