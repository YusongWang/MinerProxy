package controllers

import (
	"log"
	"miner_proxy/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Logger(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()
	//ws.WriteMessage(websocket.TextMessage, []byte("Hello world"))
	t, err := tail.TailFile("./logs/MinerProxy.log", tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}})
	if err != nil {
		utils.Logger.Info("tail file failed, err: " + err.Error())
		return
	}

	for line := range t.Lines {
		log.Println(line.Text)
		ws.WriteMessage(websocket.TextMessage, []byte(line.Text))
	}
}
