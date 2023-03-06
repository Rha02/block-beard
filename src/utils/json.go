package utils

import "encoding/json"

func JsonStatus(msg string) []byte {
	m, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: msg,
	})
	return m
}
