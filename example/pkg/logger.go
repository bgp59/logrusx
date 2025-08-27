package pkg

import (
	"github.com/bgp59/logrusx"
)

// Create the root logger:
var RootLogger = logrusx.NewCollectableLogger()

// Use init to add the module's root path (2 dirs up from here) to the list of
// prefixes to be stripped when logging the caller's source file.
func init() {
	RootLogger.AddCallerSrcPathPrefix(2)
}
