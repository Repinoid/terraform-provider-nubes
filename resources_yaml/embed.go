package resources_yaml

import "embed"

// Files embeds unified service YAML specs for provider runtime lookup.
//
//go:embed *.yaml
var Files embed.FS
