package diagnostics

import "github.com/davidsbond/kollect/internal/flag"

// Flags contains all command-line flags that can be used to configure diagnostics.
var Flags = flag.Flags{
	&flag.Boolean{
		Name:        "diagnostics",
		Usage:       "Enables pprof endpoints at /pprof/* for diagnostics",
		EnvVar:      "ENABLE_DIAGNOSTICS",
		Destination: &enableDiagnostics,
		Hidden:      true,
	},
}

var enableDiagnostics bool
