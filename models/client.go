package models

import "encoding/json"

//Client 定义
type Client struct {
	ID            string `json:"token"`
	Type          string `json:"type,omitempty"`
	CurrentLayout string `json:"currentlayout,omitempty"`
	IP            string `json:"ip,omitempty"`
	MAC           string `json:"mac,omitempty"`
}

const (
	//CLIENT_TYPE_CONTROL makes fff
	ClientTypeControl = "control"
	ClientTypeDisplay = "display"
)

type Clients struct {
	Map map[string]Client
}

//IsExistClient 判断clientMap中是否存在该key，参数key,返回值bool
func (c *Clients) IsExistClient(id string) bool {
	_, ok := c.Map[id]
	return ok
}

func (c *Clients) IsExistDisplay() bool {
	for _, v := range c.Map {
		if v.Type == ClientTypeDisplay {
			return true
		}
	}
	return false
}

//AddClient clientMap中增加client
func (c *Clients) AddClient(s string) string {
	i := &Client{}
	json.Unmarshal([]byte(s), i)
	c.Map[i.ID] = Client{
		i.ID,
		i.Type,
		i.CurrentLayout,
		i.IP,
		i.MAC,
	}

	return i.CurrentLayout
}

//RemoveClient clientMap中移除client
func (c *Clients) RemoveClient(k string) {
	delete(c.Map, k)
}
