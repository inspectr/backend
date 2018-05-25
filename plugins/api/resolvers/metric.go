package inspectr_resolvers

import (
	"encoding/json"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/jinzhu/gorm"
)

// Metric
type Metric struct {
	// Interval
	Interval int32 `json:"interval"`
	// StartsAt
	StartsAt time.Time `json:"startsAt"`
	// EventMetadata
	Size int32 `json:"size"`
}

// TrailResolver resolver for Trail
type MetricResolver struct {
	Metric
	DB *gorm.DB
}

// Interval
func (r *MetricResolver) Interval() int32 {
	return int32(r.Metric.Interval)
}

// StartsAt
func (r *MetricResolver) StartsAt() graphql.Time {
	return graphql.Time{Time: r.Metric.StartsAt}
}

// EventMetadata
func (r *MetricResolver) Size() int32 {
	return int32(r.Metric.Size)
}

func (r *MetricResolver) MarshalJSON() ([]byte, error) {
	return json.Marshal(&r.Metric)
}

func (r *MetricResolver) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.Metric)
}
