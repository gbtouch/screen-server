package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
)

const (
	DBUrl  = "mongodb://192.168.19.66:27017"
	DBName = "appDatabase"
	//DBUrl  = "mongodb://127.0.0.1:27017"
)

//Open makes screen_server open
func Open() {
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

	if err != nil {
		log.Fatal(err)
	}

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
