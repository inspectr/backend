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

type Trail struct {
	// Timestamp
	Timestamp int64 `json:"timestamp"`
	// Event
	Event string `json:"event"`
	// EventMetadata
	EventMetadata interface{} `json:"eventMetadata"`
	// Actor
	Actor string `json:"actor"`
	// ActorMetadata
	ActorMetadata interface{} `json:"actorMetadata"`
	// Target
	Target string `json:"target"`
	// TargetMetadata
	TargetMetadata interface{} `json:"targetMetadata"`
	// Origin
	Origin string `json:"origin"`
	// OriginMetadata
	OriginMetadata interface{} `json:"originMetadata"`

	// MessageID
	MessageID string
}
