package models

import "gopkg.in/mgo.v2/bson"

//Grid make a grid
type Grid struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Left     int           `json:"left"`
	Top      int           `json:"top"`
	ZIndex   int           `json:"zindex"`
	Width    int           `json:"width"`
	Height   int           `json:"height"`
	Resource *Resource     `json:"resource" bson:"resource"`
}
