package models

type GetWsDataChan struct {
	Send map[string]interface{} `json:"send"`
}
