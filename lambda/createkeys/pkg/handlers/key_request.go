package handlers

import (
	"encoding/json"
	"errors"
)

type keyRequestEvent struct {
	PublicKey string `json:"public_key"`
}

func (k keyRequestEvent) String() string {
	b, _ := json.Marshal(k)
	return string(b)
}

var errEmptyPublicKey = errors.New("PublicKey attribute is empty or missing")

func Unpack(body string) (interface{}, error) {
	var event keyRequestEvent
	err := json.Unmarshal([]byte(body), &event)
	if err != nil {
		return nil, err
	}
	if len(event.PublicKey) == 0 {
		return nil, errEmptyPublicKey
	}

	return &event, nil
}
