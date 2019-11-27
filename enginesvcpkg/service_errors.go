package enginesvcpkg

import "errors"

var (
	ErrEngineNotExistInRegistry = errors.New("Engine not exist in registry")
	ErrEnginePreExistInRegistry = errors.New("Engine pre-exists in registry")
	ErrInconsistentEngineIDs	= errors.New("Inconsistent Engine IDs")
)

