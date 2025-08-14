package app

import (
	"context"
	"net/http"

	"sisupass.com/sisupass/internal/types"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *Application) ContextSetUser(r *http.Request, user *types.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *Application) ContextGetUser(r *http.Request) *types.User {
	user, ok := r.Context().Value(userContextKey).(*types.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
