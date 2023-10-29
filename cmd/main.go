package main

import (
	"qqbot/sever/baidusever"
	"qqbot/sever/messagesever"
	"qqbot/sever/wssever"
)

func main() {
	baidusever.InitEnv()
	dataChan := wssever.GetWsResMessage()
	for data := range dataChan {
		go messagesever.SendMsg(data)
	}

}
