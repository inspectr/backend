package inspectr_resolvers

import (
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
