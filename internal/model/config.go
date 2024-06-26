package model

import "time"

type Config struct {
	Log           *Log         `yaml:"log"`
	RootDir       string       `yaml:"root-dir"`
	Authorization string       `yaml:"authorization"`
	SelfAPI       *SelfAPI     `yaml:"self-api"`
	TelegramBot   *TelegramBot `yaml:"telegram-bot"`
}

type Log struct {
	ForceNew bool          `yaml:"force-new,omitempty"`
	Level    string        `yaml:"level,omitempty"`
	Aging    time.Duration `yaml:"aging,omitempty"`
	Colorful bool          `yaml:"colorful,omitempty"`
}

type SelfAPI struct {
	RootDir     string `yaml:"root-dir,omitempty"`
	RealRootDir string `yaml:"real-root-dir,omitempty"`
}

type TelegramBot struct {
	Debug    bool   `yaml:"debug,omitempty"`
	Endpoint string `yaml:"endpoint,omitempty"`
	Token    string `yaml:"token,omitempty"`
}
