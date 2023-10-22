package di

import (
	"context"
)

// Get the dependency from the dependency container using the key
func Get(ctx context.Context, key string) any {
	ctn, ok := ctx.Value(containerKey).(*container)
	if !ok {
		panic("container does not exist on context")
	}

	return ctn.Get(key)
}
