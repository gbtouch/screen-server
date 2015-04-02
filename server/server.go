package server

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"screen-server/log"

	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/yaml.v2"

	l4g "github.com/alecthomas/log4go"
)

type Config struct {
	DB map[string][]string `yaml:"db,omitempty"`
}

var Settings Config

func check(e error) {
	if e != nil {
		l4g.Error("Server err: %v", e)
	}
}

func loadConfig() {
	filename, _ := filepath.Abs("./config.yaml")

	yamlFile, err := ioutil.ReadFile(filename)
	check(err)

	err = yaml.Unmarshal(yamlFile, &Settings)
	check(err)
}

//Open makes screen_server open
func Open() {
	logger.Init()
	l4g.Debug("Start Screen Server")
	l4g.Info("loading local config")
	loadConfig()
	l4g.Info("loaded")
	l4g.Info("--------------------")

	l4g.Info("connect mongo db")
	Mongo.InitDB()
	l4g.Info("connected")
	l4g.Info("--------------------")

	l4g.Info("loading resource&layouts")
	Resources.Load(Mongo.DB)
	Layouts.Load(Mongo.DB)
	l4g.Info("loaded")
	l4g.Info("--------------------")

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		&rest.Route{"GET", "/resource", getResourceHandler},
		&rest.Route{"POST", "/resource", setResourceHandler},
		&rest.Route{"GET", "/layouts", getLayoutsHandler},
		&rest.Route{"POST", "/layout", setLayoutHandler},
		&rest.Route{"POST", "/layout/current", setCurrentLayoutHandler},
		&rest.Route{"PATCH", "/layout/current", updateCurrentLayoutHandler},
		&rest.Route{"POST", "/layout/:id/resource", updateLayoutResourceHandler},
		&rest.Route{"POST", "/error", notifyErrorHandler},
		&rest.Route{"POST", "/control/heartbeat", setHeartbeatHandler},
		&rest.Route{"POST", "/display/heartbeat", setHeartbeatHandler},
	)

	check(err)

	api.SetApp(router)

	//Mock()

	server = newSocketServer()

	// ticker := time.NewTicker(time.Second * 10)
	//
	// go func() {
	// 	for t := range ticker.C {
	// 		logger.Log.Info("time:", t)
	// 		logger.Log.Info("clients:", clients.Map)
	// 	}
	// }()

	http.Handle("/", api.MakeHandler())
	http.Handle("/socket.io/", server)
	err = http.ListenAndServe(":8080", nil)
	check(err)

	l4g.Info("Start Screen Server")
}
