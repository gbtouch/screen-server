package server

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gbtouch/screen-server/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DB         map[string][]string `yaml:"db,omitempty"`
	ServiceUrl []string            `yaml:"configserivce"`
}

var Settings Config

func check(e error) {
	if e != nil {
		logger.Log.Error("Server err: %v", e)
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
	logger.Log.Debug("Start Screen Server")
	logger.Log.Info("loading local config")
	loadConfig()
	logger.Log.Info("loaded")
	logger.Log.Info("--------------------")

	logger.Log.Info("connect mongo db")
	Mongo.InitDB()
	logger.Log.Info("connected")
	logger.Log.Info("--------------------")

	logger.Log.Info("loading resource&layouts")
	if Mongo.DB != nil {
		Resources.Load(Mongo.DB)
		Layouts.Load(Mongo.DB)
		logger.Log.Info("loaded")
	} else {
		logger.Log.Error("connect to mongodb failed")
	}
	logger.Log.Info("--------------------")

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

	ticker := time.NewTicker(time.Second * 30)

	go func() {
		for t := range ticker.C {
			logger.Log.Debug("time:", t)
			logger.Log.Debug("clients:", clients.Map)
		}
	}()

	http.Handle("/", api.MakeHandler())
	http.Handle("/socket.io/", server)
	err = http.ListenAndServe(":8080", nil)
	check(err)

	logger.Log.Info("Start Screen Server")
}
