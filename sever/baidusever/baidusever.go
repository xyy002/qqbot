package baidusever

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"qqbot/models"
	"strings"
)

var Token = ""

func init() {
	//InitEnv()
	Token = GetAssessToken()
}
func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// GetAssessToken 获取百度AI的access_token token有效期为30天 没有写自动刷新 到期了自己手动刷新
func GetAssessToken() string {
	InitEnv()
	id := os.Getenv("BAIDU_CLIENT_ID")
	secret := os.Getenv("BAIDU_CLIENT_SECRET")
	if id == "" || secret == "" {
		panic("请在.env文件中配置百度AI的client_id和client_secret")
	}
	//log.Print(id, secret)
	url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", id, secret)
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	//fmt.Println(string(body))
	var result map[string]interface{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		log.Println("Error:", err)
		return ""
	}
	t := result["access_token"].(string)

	return t
}

func GetMsg(token string, message []models.BdMessage) (string, error) {
	url := fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions?access_token=%s", token)
	method := "POST"

	payloadData := models.BdMsgData{
		Messages: message,
	}

	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return "", err
	}

	payload := bytes.NewReader(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to do request: %v", err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)
	// 如果stream为false，直接读取响应的主体并返回
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return "", err
	}
	var result map[string]interface{}
	if err = json.Unmarshal([]byte(body), &result); err != nil {
		log.Println(err)
		return "", err
	}
	log.Println(result["results"])
	return result["result"].(string), nil

}
