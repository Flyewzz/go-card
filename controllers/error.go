package controllers

import (
	"encoding/json"
	"net/http"
)

type ErrMsg struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func HandleError(w *http.ResponseWriter, r *http.Request, code int, text string) {
	msg := ErrMsg{
		Code: code,
		Text: text,
	}
	data, _ := json.Marshal(msg)
	(*w).WriteHeader(msg.Code)
	(*w).Write(data)
}
