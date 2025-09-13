package types

import (
	shared "govel/types/shared"
)

/**
 * ShutdownCallback is a type alias for the shared version to maintain backward compatibility.
 * The actual type definition is now in the shared package to avoid circular imports.
 */
type ShutdownCallback = shared.ShutdownCallback
