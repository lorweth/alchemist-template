package iam

import (
	"context"
)

const (
	userProfileCtxKey = "user_profile"
)

func SetInCtx(ctx context.Context, profile UserProfile) context.Context {
	return context.WithValue(ctx, userProfileCtxKey, profile)
}

func FromCtx(ctx context.Context) UserProfile {
	p, ok := ctx.Value(userProfileCtxKey).(UserProfile)
	if !ok {
		return UserProfile{}
	}

	return p
}
