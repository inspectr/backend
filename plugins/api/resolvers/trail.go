package inspectr_resolvers

import (
	"encoding/json"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// Trail
type Trail struct {
	Model `json:",inline"`
	// Timestamp
	Timestamp int64 `json:"timestamp" gorm:"type:integer"`
	// Event
	Event string `json:"event" gorm:"type:varchar(100)"`
	// EventMetadata
	EventMetadata postgres.Jsonb `json:"eventMetadata" gorm:"type:jsonb;"`
	// Actor
	Actor string `json:"actor"`
	// ActorMetadata
	ActorMetadata postgres.Jsonb `json:"actorMetadata" gorm:"type:jsonb;"`
	// Target
	Target string `json:"target"`
	// TargetMetadata
	TargetMetadata postgres.Jsonb `json:"targetMetadata" gorm:"type:jsonb;"`
	// Origin
	Origin string `json:"origin"`
	// OriginMetadata
	OriginMetadata postgres.Jsonb `json:"originMetadata" gorm:"type:jsonb;"`
}

// TrailResolver resolver for Trail
type TrailResolver struct {
	Trail
	DB *gorm.DB
}

// ID
func (r *TrailResolver) ID() graphql.ID {
	return graphql.ID(r.Trail.Model.ID.String())
}

// Timestamp
func (r *TrailResolver) Timestamp() int64 {
	return r.Trail.Timestamp
}

// Event
func (r *TrailResolver) Event() string {
	return r.Trail.Event
}

// EventMetadata
func (r *TrailResolver) EventMetadata() JSON {
	return JSON{r.Trail.EventMetadata.RawMessage}
}

// Actor
func (r *TrailResolver) Actor() string {
	return r.Trail.Actor
}

// ActorMetadata
func (r *TrailResolver) ActorMetadata() JSON {
	return JSON{r.Trail.ActorMetadata.RawMessage}
}

// Target
func (r *TrailResolver) Target() string {
	return r.Trail.Target
}

// TargetMetadata
func (r *TrailResolver) TargetMetadata() JSON {
	return JSON{r.Trail.TargetMetadata.RawMessage}
}

// Origin
func (r *TrailResolver) Origin() string {
	return r.Trail.Origin
}

// OriginMetadata
func (r *TrailResolver) OriginMetadata() JSON {
	return JSON{r.Trail.OriginMetadata.RawMessage}
}

// Created
func (r *TrailResolver) Created() graphql.Time {
	tm := time.Unix(r.Trail.Timestamp, 0)
	return graphql.Time{Time: tm}
}

func (r *TrailResolver) MarshalJSON() ([]byte, error) {
	return json.Marshal(&r.Trail)
}

func (r *TrailResolver) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.Trail)
}
