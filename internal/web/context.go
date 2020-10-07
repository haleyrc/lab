package web

import (
	"context"
)

type key int

const (
	userKey key = iota
)

func getCurrentUserFromContext(ctx context.Context) string {
	tmp := ctx.Value(userKey)
	if tmp == nil {
		return ""
	}
	return tmp.(string)
}

func setCurrentUserOnContext(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, userKey, uid)
}
