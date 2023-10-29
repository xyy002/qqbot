package wssever

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"qqbot/models"
	"sync"
	"time"
)

func GetWsResMessage() chan models.GetWsDataChan {
	var wg sync.WaitGroup

	u := url.URL{Scheme: "wss", Host: "api.sgroup.qq.com", Path: "/websocket"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	var lastS int
	ticker := time.NewTicker(10 * time.Second)

	messageChan := make(chan []byte)
	done := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				done <- true
				return
			}
			messageChan <- message
		}
	}()

	data := make(chan models.GetWsDataChan, 1) // Add buffer size to the channel

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case message := <-messageChan:
				var serverMsg map[string]interface{}
				if err := json.Unmarshal(message, &serverMsg); err != nil {
					log.Println("read:", err)
					continue
				}

				if serverMsg["s"] != nil && serverMsg["op"] != 11 {
					lastS = int(serverMsg["s"].(float64))
				}

				if serverMsg["t"] != nil && serverMsg["t"].(string) == "AT_MESSAGE_CREATE" {
					sendData := models.GetWsDataChan{
						Send: map[string]interface{}{
							"id":             serverMsg["id"].(string),
							"authorID":       serverMsg["d"].(map[string]interface{})["author"].(map[string]interface{})["id"].(string),
							"authorUsername": serverMsg["d"].(map[string]interface{})["author"].(map[string]interface{})["username"].(string),
							"content":        serverMsg["d"].(map[string]interface{})["content"].(string),
							"channel_id":     serverMsg["d"].(map[string]interface{})["channel_id"].(string),
						},
					}

					data <- sendData
					//log.Print(serverMsg["d"].(map[string]interface{})["content"].(string))
				}
				//log.Printf("s:%d\n", lastS)
				//fmt.Printf("recv: %s\n", message)

			case <-done:
				close(data)
				return
			}
		}
	}()
	id := os.Getenv("QQ_BOT_APPID")
	token := os.Getenv("QQ_BOT_TOKEN")
	if id == "" || token == "" {
		panic("请在.env文件中配置QQ机器人的appid和token")
	}
	auth := fmt.Sprintf("Bot %s.%s", id, token)
	msg := map[string]interface{}{
		"op": 2,
		"d": map[string]interface{}{
			"token":   auth,
			"intents": 1073741824,
			"shard":   []int{0, 1},
			"properties": map[string]string{
				"$os":      "linux",
				"$browser": "my_library",
				"$device":  "my_library",
			},
		},
	}

	if err := c.WriteJSON(msg); err != nil {
		log.Println("write:", err)
		return nil
	}

	// Start a new goroutine to send heartbeat messages
	go func() {
		for {
			select {
			case <-ticker.C:
				heartbeatMsg := map[string]interface{}{
					"op": 1,
					"d":  lastS,
				}
				if err := c.WriteJSON(heartbeatMsg); err != nil {
					log.Println("write:", err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Start a new goroutine to wait for all other goroutines to complete and close the connection
	go func() {
		// Wait for all goroutines to complete.
		wg.Wait()

		// Now it's safe to close the connection.
		c.Close()

		// Stop the ticker
		ticker.Stop()
	}()

	// Return the data channel immediately
	return data
}
