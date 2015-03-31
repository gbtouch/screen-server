package server

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"screen-server/log"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DB map[string][]string `yaml:"db,omitempty"`
}

var Settings Config

func check(e error) {
	if e != nil {
		logger.Log.Error("Server err:%v", e)
		return
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

	logger.Log.Info("loading local config")
	loadConfig()
	logger.Log.Info("loaded")
	logger.Log.Info("--------------------")

	logger.Log.Info("connect mongo db")
	Mongo.InitDB()
	logger.Log.Info("connected")
	logger.Log.Info("--------------------")

	logger.Log.Info("loading resource&layouts")
	Resources.Load(Mongo.DB)
	Layouts.Load(Mongo.DB)
	logger.Log.Info("loaded")
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
		&rest.Route{"POST", "/error", notifyErrorHandler})

	check(err)

	api.SetApp(router)

	//Mock()

	server = newSocketServer()

	ticker := time.NewTicker(time.Second * 10)

	go func() {
		for t := range ticker.C {
			logger.Log.Info("time:", t)
			logger.Log.Info("clients:", clients.Map)
		}
	}()

	http.Handle("/", api.MakeHandler())
	http.Handle("/socket.io/", server)
	err = http.ListenAndServe(":8080", nil)
	check(err)

	logger.Log.Info("Start screen server")
}
