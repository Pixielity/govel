package types

import (
	"context"
)

/**
 * ShutdownCallback represents a function that can be registered to run during shutdown.
 *
 * @param ctx context.Context The context for the shutdown callback
 * @return error Any error that occurred during the callback execution
 */
type ShutdownCallback func(ctx context.Context) error