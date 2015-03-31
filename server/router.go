package server

import (
	"net/http"
	"screen-server/log"
	"screen-server/models"

	"github.com/ant0ine/go-json-rest/rest"
)

var (
	//Resources make server's resource list
	Resources = models.ResourceMap{}

	//Layouts make server's layout list
	Layouts = models.Layouts{}

	ChangedLayout = models.ChangedLayout{}

	CurrentLayout = models.Layout{}

	UpdateLayout = models.UpdateLayout{}

	Mongo DBImpl

	tokens = make(map[string]models.Token)

	clients = models.Clients{
		make(map[string]models.Client),
	}
)

//Mock make mock data
// func Mock() {
// 	Resources.Mock()
// 	Layouts.Mock()
// }

//getLayoutsHandler makes T3
//last update 2015-3-26
//by tommy
func getLayoutsHandler(w rest.ResponseWriter, r *rest.Request) {
	logger.Log.Debug("GET /T3-[Header]", r.Header)
	w.WriteJson(Layouts.Store)
}

//getResourceHandler makes T4
//last update 2015-3-26
//by tommy
func getResourceHandler(w rest.ResponseWriter, r *rest.Request) {
	logger.Log.Debug("GET /T4-[Header]", r.Header)
	w.WriteJson(Resources.Store)
}

//setCurrentLayoutHandler makes T5
//last update 2015-3-27
//by tommy
func setCurrentLayoutHandler(w rest.ResponseWriter, r *rest.Request) {
	logger.Log.Debug("POST /T5-[Header]", r.Header)
	c := models.ChangedLayout{}
	err := r.DecodeJsonPayload(&c)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !Layouts.IsExistLayout(c.ID) {
		rest.Error(w, "server is currently no layout", http.StatusNotFound)
		return
	}

	if !clients.IsExistDisplay() {
		rest.Error(w, "display has not been connected to the server", http.StatusMethodNotAllowed)
		return
	}

	ChangedLayout = c

	server.BroadcastTo("warroom", "LayoutChanged", &c)

	w.WriteJson(&c)
}

//setLayoutHandler makes T6
//last update 2015-3-26
//by tommy
func setLayoutHandler(w rest.ResponseWriter, r *rest.Request) {
	logger.Log.Debug("POST /T6-[Header]", r.Header)
	t := models.Token{}

	err := r.DecodeJsonPayload(&t)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, ok := tokens[t.Token]

	if ok {
		rest.Error(w, "duplicate submission", http.StatusConflict)
		return
	}

	tokens[t.Token] = t

	Layouts.Load(Mongo.DB)

	server.BroadcastTo("warroom", "LayoutUpdated", &t)

	delete(tokens, t.Token)
}

//setResourceHandler makes T7
//last update 2015-3-26
//by tommy
func setResourceHandler(w rest.ResponseWriter, r *rest.Request) {
	logger.Log.Debug("POST /T7-[Header]", r.Header)
	t := models.Token{}

	err := r.DecodeJsonPayload(&t)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, ok := tokens[t.Token]

	if ok {
		rest.Error(w, "duplicate submission", http.StatusConflict)
		return
	}

	tokens[t.Token] = t

	Resources.Load(Mongo.DB)

	server.BroadcastTo("warroom", "ResourceUpdated", &t)

	delete(tokens, t.Token)
	w.WriteJson(t.Token)
}

//notifyErrorHandler makes T8
//last update 2015-3-26
//by tommy
func notifyErrorHandler(w rest.ResponseWriter, r *rest.Request) {
	logger.Log.Debug("POST /T8-[Header]", r.Header)
	error := models.Error{}
	err := r.DecodeJsonPayload(&error)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	server.BroadcastTo("warroom", "ErrorNotify", &error)
	w.WriteJson(&error)
}

//updateLayoutGridsHandler makes T9
//last update 2015-3-27
//by tommy
func updateLayoutResourceHandler(w rest.ResponseWriter, r *rest.Request) {
	logger.Log.Debug("PATCH /T9-[Header]", r.Header)
	id := r.PathParam("id")

	if CurrentLayout.ID != id {
		rest.Error(w, "layout not found by id", http.StatusNotFound)
		return
	}

	err := r.DecodeJsonPayload(&UpdateLayout)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	UpdateLayout.ID = id

	server.BroadcastTo("warroom", "LayoutResourceUpdated", &UpdateLayout)
	w.WriteJson(&UpdateLayout)
}

//updateCurrentLayoutHandler makes T10
//last update 2015-3-27
//by tommy
func updateCurrentLayoutHandler(w rest.ResponseWriter, r *rest.Request) {
	logger.Log.Debug("PATCH /T10-[Header]", r.Header)
	d := models.Layout{}

	err := r.DecodeJsonPayload(&d)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if CurrentLayout.ID != d.ID {
		rest.Error(w, "current layout not found", http.StatusNotFound)
	}

	CurrentLayout = d
	server.BroadcastTo("warroom", "CurrentLayoutUpdated", &CurrentLayout)
	w.WriteJson(&CurrentLayout)
}
