package contextutil

import (
	"context"
	"github.com/ruhollahh/paperback/internal/app/domain"
)

type contextKey string

const userContextKey = contextKey("user")
const nonceContextKey = contextKey("nonce")

func ContextSetUser(c context.Context, user *domain.User) context.Context {
	return context.WithValue(c, userContextKey, user)
}

func ContextGetUser(c context.Context) *domain.User {
	user, ok := c.Value(userContextKey).(*domain.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}

func ContextSetNonce(c context.Context, nonce string) context.Context {
	return context.WithValue(c, nonceContextKey, nonce)
}

func ContextGetNonce(c context.Context) string {
	nonce, ok := c.Value(nonceContextKey).(string)
	if !ok {
		panic("missing nonce value in request context")
	}

	return nonce
}
