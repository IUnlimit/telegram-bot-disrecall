package bot

import (
	global "github.com/IUnlimit/telegram-bot-disrecall/internal"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/db"
	log "github.com/sirupsen/logrus"
)

func UserAuth(chatID int64, bot *BasicTGBot) bool {
	authorization := global.Config.Authorization
	if authorization == "" {
		return true
	}
	exist, err := db.ExistUser(chatID)
	if err != nil {
		log.Error(err)
		bot.SendChatMessage("鉴权失败, 熔断所有访问!", chatID)
		return false
	}
	if exist {
		return true
	}

	channel := make(chan string)
	InputChanMap[chatID] = channel
	bot.SendChatMessage("限制访问已开启, 请输入验证密匙", chatID)

	inputToken := <-channel
	delete(InputChanMap, chatID)
	if inputToken != authorization {
		bot.SendChatMessage("密匙错误, 请重新输入!", chatID)
		return false
	}
	err = db.InsertUser(chatID)
	if err != nil {
		bot.SendChatMessage("未知错误, 请重新验证", chatID)
		return false
	}
	bot.SendChatMessage("验证成功, 输入 /help 查看帮助", chatID)
	return false
}
