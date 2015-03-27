package models

import (
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type resourceType int

const (
	media resourceType = iota
	image
	clock
)

//Resource defines the basic information resource
type Resource struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	Type resourceType  `json:"stype" bson:"stype"`
	URL  string        `json:"url" bson:"url"`
	Fix  bool          `json:"fix" bson:"fix"`
}

//ResourceMap make server's resource list
type ResourceMap struct {
	sync.RWMutex
	Store   map[string]Resource
	Created string
}

//Load 方法实现从MongoDB读取资源集合(collection name: resources)填充全局变量ResourceMap
func (r *ResourceMap) Load(d *mgo.Database) {
	r.Lock()
	var l []Resource
	r.Store = map[string]Resource{}
	d.C("resources").Find(nil).All(&l)
	for i := range l {
		r.Store[l[i].ID.Hex()] = l[i]
	}
	r.Created = time.Now().Format("2006-01-02 15:04:05")
	r.Unlock()
	//log.Println("resource create time:", r.Created)
}

//Mock makes some demo data
// func (r *ResourceMap) Mock() {
// 	guid := [3]bson.ObjectId{}
// 	for i := 0; i < 3; i++ {
// 		guid[i] = bson.NewObjectId()
// 	}
//
// 	r.Lock()
// 	r.Store = map[string]Resource{}
// 	r.Store[guid[0].Hex()] = Resource{
// 		guid[0],
// 		image,
// 		"http://192.168.1.1/images/001.png",
// 		true}
// 	r.Store[guid[1].Hex()] = Resource{
// 		guid[1],
// 		media,
// 		"rtsp://192.168.16.140:5554/stream.smp?address=192.168.16.181&channel=0",
// 		false}
// 	r.Store[guid[2].Hex()] = Resource{
// 		guid[2],
// 		media,
// 		"rtsp://192.168.16.140:5554/stream.smp?address=192.168.16.182&channel=0",
// 		false}
//
// 	r.Unlock()
// }
