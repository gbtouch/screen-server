package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	DB map[string][]string `yaml:"db,omitempty"`
}

var Config Configuration

func check(e error) {
	if e != nil {
		log.Println(e)
	}
}

func loadConfig() {
	filename, _ := filepath.Abs("./config.yaml")

	yamlFile, err := ioutil.ReadFile(filename)
	check(err)

	err = yaml.Unmarshal(yamlFile, &Config)
	check(err)
}

//Open makes screen_server open
func Open() {
	loadConfig()

	Mongo.InitDB()

	Resources.Load(Mongo.DB)
	Layouts.Load(Mongo.DB)

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

	ticker := time.NewTicker(time.Second * 3)

	go func() {
		for t := range ticker.C {
			fmt.Println("Loop", t)
			fmt.Println("Clients:", clients.Map)
		}
	}()

	http.Handle("/", api.MakeHandler())
	http.Handle("/socket.io/", server)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
