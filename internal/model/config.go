package model

import "time"

type Config struct {
	Log         *Log         `yaml:"log"`
	API         *API         `yaml:"api"`
	RootDir     string       `yaml:"root-dir"`
	TelegramBot *TelegramBot `yaml:"telegram-bot"`
}

type Log struct {
	ForceNew bool          `yaml:"force-new,omitempty"`
	Level    string        `yaml:"level,omitempty"`
	Aging    time.Duration `yaml:"aging,omitempty"`
	Colorful bool          `yaml:"colorful,omitempty"`
}

type API struct {
	Host string `yaml:"host,omitempty"`
	Port int64  `yaml:"port,omitempty"`
}

type TelegramBot struct {
	Debug    bool   `yaml:"debug,omitempty"`
	Endpoint string `yaml:"endpoint,omitempty"`
	Token    string `yaml:"token,omitempty"`
}
