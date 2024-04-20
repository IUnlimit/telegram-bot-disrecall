package main

import (
	"github.com/IUnlimit/telegram-bot-disrecall/cmd/disrecall"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/conf"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/logger"
)

func main() {
	conf.Init()
	logger.Init()
	disrecall.Listen()
}