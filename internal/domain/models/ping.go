package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type PingData struct {
	Id         uuid.UUID           `json:"id"`
	CreatedAt  time.Time           `json:"created_at"`
	IP         string              `json:"ip"`
	RequestURI string              `json:"request_uri"`
	UserAgent  string              `json:"user_agent"`
	Headers    map[string][]string `json:"headers"`
	Client     string              `json:"client"`
	Device     string              `json:"device"`
	OS         string              `json:"os"`
	Lang       []string            `json:"lang"`
	Country    string              `json:"country"`
}

func (d PingData) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

type CreatePingData struct {
	IP         string              `json:"ip"`
	RequestURI string              `json:"request_uri"`
	UserAgent  string              `json:"user_agent"`
	Headers    map[string][]string `json:"headers"`
}
