package common

import (
	"strconv"
	"time"
)

type Metadata struct {
	Message string `json:"message"`
	Time    string `json:"time"`
	Ecode   string `json:"ecode"`
}

// 定制统一的返回体内容
func NewResMessage(ecode string, msg string, payload interface{}) map[string]interface{} {
	body := make(map[string]interface{})
	//metadata := make(map[string]Metadata)
	//metadata["request_time"] = time.Now().UnixNano() / 1e6

	if msg == "" {
		msg = "Successfully."
	}
	metadata := Metadata{
		Message: msg,
		Time:    strconv.FormatInt(time.Now().UnixNano()/1e6, 10),
		Ecode:   ecode,
	}
	body["metadata"] = metadata
	if payload != "" && payload != nil {
		body["payload"] = payload
	}

	return body
}

func ResBody(ecode string, msg string, payload interface{}) map[string]interface{} {
	return NewResMessage(ecode, msg, payload)
}
