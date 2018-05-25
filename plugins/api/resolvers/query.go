package inspectr_resolvers

import (
	"context"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
)

// User Retrieve single user by ID
func (r *Resolver) User(args *struct {
	ID *graphql.ID
}) *UserResolver {
	return nil
}

// Users Retrieve all users
func (r *Resolver) Users() []*UserResolver {
	return nil
}

// Trails
func (r *Resolver) Trails(ctx context.Context) ([]*TrailResolver, error) {
	var rows []Trail
	var results []*TrailResolver

	r.DB.Order("created_at desc").Find(&rows)

	for _, trail := range rows {
		results = append(results, &TrailResolver{DB: r.DB, Trail: trail})
	}

	return results, nil
}

// Metrics
func (r *Resolver) Metrics(ctx context.Context, args *struct {
	StartsAt *graphql.Time
	EndsAt   *graphql.Time
	Interval *int32
}) ([]*MetricResolver, error) {
	//var rows []Trail
	var results []*MetricResolver

	type IntervalMetric struct {
		Count    int32     `json:"count"`
		Interval time.Time `json:"interval"`
	}

	var intervalMetrics []IntervalMetric
	_intervalMetrics := make(map[int64]Metric)

	r.DB.Debug().Raw(`SELECT COUNT(*) count, 
	to_timestamp(floor((extract('epoch' from created_at) / $1 )) * $1) 
	as interval
	FROM trails
	WHERE created_at >= $2 AND created_at <= $3
	GROUP BY interval`, args.Interval, args.StartsAt.Time, args.EndsAt.Time).Scan(&intervalMetrics)

	for _, m := range intervalMetrics {
		_intervalMetrics[m.Interval.Unix()] = Metric{
			Interval: *args.Interval,
			StartsAt: m.Interval,
			Size:     m.Count,
		}

	}

	var metrics []Metric

	t := args.StartsAt.Time
	endsAt := args.EndsAt.Time

	for t.Before(endsAt) {
		m := Metric{
			Interval: *args.Interval,
			StartsAt: t,
			Size:     0,
		}

		if val, ok := _intervalMetrics[t.Unix()]; ok {
			m = val
		}

		metrics = append(metrics, m)

		t = t.Add(time.Duration(*args.Interval) * time.Second)
	}

	for _, m := range metrics {
		results = append(results, &MetricResolver{DB: r.DB, Metric: m})
	}

	return results, nil
}
