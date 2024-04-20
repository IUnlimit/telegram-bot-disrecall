package configs

import "embed"

//go:embed config.yml
var Config embed.FS

var ConfigFileName = "config.yml"
