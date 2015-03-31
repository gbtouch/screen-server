package server

import (
	"encoding/json"
	"screen-server/models"

	"github.com/gbtouch/go-socket.io"
)

//clients is connected to the service side of the Socket map
var server *socketio.Server

func newSocketServer() *socketio.Server {
	socket, err := socketio.NewServer()

	check(err)

	socket.On("connection", func(so socketio.Socket) {
		so.Join("warroom")
		token, _ := json.Marshal(map[string]string{"token": so.Id()})
		so.Emit("IntializeRequest", string(token))

		so.On("IntializeResponse", func(msg string) {
			if !clients.IsExistClient(so.Id()) {
				so.Emit("IntializeCompleted", clients.AddClient(msg))
			} else {
				so.Emit("IntializeCompleted", "terminal repeat")
			}
		})

		so.On("LayoutResponse", func(msg string) {
			var s models.ResponseLayout

			json.Unmarshal([]byte(msg), &s)

			if s.Action == "update" && s.Result {
				//修改Grid
				CurrentLayout.UpdateGrid(&UpdateLayout)
			}

			if s.Action == "change" && s.Result {
				//修改当前layout
				CurrentLayout = Layouts.Store[s.ID]
			}

			socket.BroadcastTo("warroom", "LayoutResponse", msg)
		})

		so.On("disconnection", func() {
			so.Leave("warroom")
			clients.RemoveClient(so.Id())
		})
	})

	return socket
}
