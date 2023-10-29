package models

type BdMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BdMsgData struct {
	//Stream   bool        `json:"stream"`
	Messages []BdMessage `json:"messages"`
}
