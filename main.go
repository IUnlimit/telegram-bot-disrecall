package main

import (
	"github.com/IUnlimit/telegram-bot-disrecall/cmd/disrecall"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/cache"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/conf"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/db"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/logger"
)

func main() {
	conf.Init()
	logger.Init()
	db.Init()
	cache.Init()
	disrecall.Listen()
}
