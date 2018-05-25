package inspectr_resolvers

import (
	"context"

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
