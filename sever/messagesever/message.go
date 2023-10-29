package messagesever

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"qqbot/models"
	"qqbot/sever/baidusever"
	"strings"
)

func SendMessage(d models.GetWsDataChan, c string) {
	channelID := d.Send["channel_id"].(string)
	if channelID == "" {
		panic("channelID为空")
	}
	url := fmt.Sprintf("https://api.sgroup.qq.com/channels/%s/messages", channelID)
	method := "POST"
	data := map[string]interface{}{
		"content": fmt.Sprintf("<@!%s> :%s", d.Send["authorID"].(string), c),
		"msg_id":  d.Send["id"].(string),
	}
	payload, _ := json.Marshal(data)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(payload))

	if err != nil {
		fmt.Println(err)
		return
	}
	id := os.Getenv("QQ_BOT_APPID")
	token := os.Getenv("QQ_BOT_TOKEN")
	if id == "" || token == "" {
		panic("请在.env文件中配置QQ机器人的appid和token")
	}
	auth := fmt.Sprintf("Bot %s.%s", id, token)
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "api.sgroup.qq.com")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	//body, err := io.ReadAll(res.Body)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(string(body))
}

func SendMsg(d models.GetWsDataChan) {
	var content = ""
	if d.Send["content"] != nil {
		content = d.Send["content"].(string)
		content = strings.Split(content, ">")[1]
		content = content[1:]
	}
	var msg = []models.BdMessage{
		{
			Role:    "user",
			Content: content,
		},
	}
	fmt.Println(baidusever.Token)
	m, _ := baidusever.GetMsg(baidusever.Token, msg)
	fmt.Println(m)
	SendMessage(d, m)
}
