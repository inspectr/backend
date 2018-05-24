package plugins

import (
	uuid "github.com/satori/go.uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	APIKey    string    `json:"APIKey" role:"secret"`
	APISecret string    `json:"APISecret" role:"secret"`
}

type HeartBeat struct {
	Tick string `json:"tick"`
}
