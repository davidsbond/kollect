package telemetry

import "github.com/davidsbond/kollect/internal/flag"

// Flags contains all command-line flags that can be used to configure telemetry.
var Flags = flag.Flags{
	&flag.String{
		Name:        "telemetry-url",
		Usage:       "URL to configure opentelemetry. See documentation for all configuration options.",
		Destination: &telemetryURL,
		EnvVar:      "TELEMETRY_URL",
		Value:       "noop://",
		Hidden:      true,
	},
}

var telemetryURL string
