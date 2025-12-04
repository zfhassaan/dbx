package notify

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SlackNotify(webhookURL, message string) error {
	payload := map[string]string{"text": message}
	data, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(data))
	return err
}
