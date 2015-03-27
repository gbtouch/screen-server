package models

import (
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Token struct {
	Token string `json:"token"`
}

//CurrentLayout make a Current Layout
type ChangedLayout struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type ResponseLayout struct {
	Result bool   `json:"result"`
	Action string `json:"action"`
	ID     string `json:"id"`
	Token  string `json:"token"`
}

type UpdateLayout struct {
	ID    string          `json:"id,omitempty"`
	Grids map[string]Grid `json:"grids" bson:"grids"`
	Token string          `json:"token" bson:"token"`
}

//Layout make a layout
type Layout struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Index int             `json:"index"`
	Grids map[string]Grid `json:"grids"`
}

type LayoutDO struct {
	ID    bson.ObjectId `bson:"_id"`
	Name  string        `bson:"name"`
	Index int           `bson:"index"`
	Grids []Grid        `bson:"grids"`
}

type Layouts struct {
	sync.RWMutex
	Store   map[string]Layout
	Created string
}

func (l *Layouts) IsExistLayout(k string) bool {
	_, ok := l.Store[k]

	return ok
}

func (r *Layouts) Load(d *mgo.Database) {
	r.Lock()
	var l []LayoutDO
	r.Store = map[string]Layout{}
	d.C("layouts").Find(nil).All(&l)
	for i := range l {
		r.Store[l[i].ID.Hex()] = createLayout(l[i])
	}
	r.Created = time.Now().Format("2006-01-02 15:04:05")
	r.Unlock()
}

func (r *Layout) UpdateGrid(g *UpdateLayout) {
	if r.ID != g.ID {
		return
	}

	//log.Println("p update------", r)

	for k, v := range g.Grids {
		r.Grids[k] = v
	}

	//log.Println("[updated]", r)
}

func createLayout(t LayoutDO) Layout {
	l := Layout{}
	l.ID = t.ID.Hex()
	l.Index = t.Index
	l.Name = t.Name
	l.Grids = createGrid(t.Grids)
	return l
}

func createGrid(g []Grid) map[string]Grid {
	r := map[string]Grid{}

	for i := range g {
		r[g[i].ID.Hex()] = g[i]
	}

	return r
}

// func (l *Layouts) Mock() {
// 	guid := [2]string{}
// 	for i := 0; i < 2; i++ {
// 		guid[i] = utils.GetGUID()
// 	}
// 	l.Lock()
// 	l.Store[guid[0]] = &Layout{
// 		guid[0],
// 		"大会模式1",
// 		0,
// 		layout0GridMock(),
// 	}
//
// 	l.Store[guid[1]] = &Layout{
// 		guid[1],
// 		"大会模式2",
// 		1,
// 		layout1GridMock(),
// 	}
// 	l.Unlock()
// }
//
// func layout0GridMock() map[string]Grid {
// 	guid := utils.GetGUID()
// 	result := map[string]Grid{}
// 	result[guid] = Grid{
// 		guid, 0, 0, 3328, 800, 0,
// 		Resource{
// 			bson.NewObjectId(),
// 			image,
// 			"http://192.168.1.1/images/001.png",
// 			true},
// 	}
// 	return result
// }
//
// func layout1GridMock() map[string]Grid {
// 	guid := [2]string{}
// 	for i := 0; i < 2; i++ {
// 		guid[i] = utils.GetGUID()
// 	}
// 	result := map[string]Grid{}
// 	result[guid[0]] = Grid{
// 		guid[0], 0, 0, 1328, 800, 0,
// 		Resource{
// 			bson.NewObjectId(),
// 			media,
// 			"rtsp://192.168.16.140:5554/stream.smp?address=192.168.16.181&channel=0",
// 			false},
// 	}
// 	result[guid[1]] = Grid{
// 		guid[1], 1328, 0, 2000, 800, 0,
// 		Resource{
// 			bson.NewObjectId(),
// 			media,
// 			"rtsp://192.168.16.140:5554/stream.smp?address=192.168.16.182&channel=0",
// 			false},
// 	}
// 	return result
// }
