package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/IUnlimit/telegram-bot-disrecall/internal/cache"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

const NULL_CALLBACK = "null"

func OnCommandCallback(callback *tgbotapi.CallbackQuery, bot *BasicTGBot) {
	log.Debugf("Callback query received from user %d with data %s", callback.From.ID, callback.Data)
	if callback.Data == NULL_CALLBACK {
		return
	}

	args := strings.Split(callback.Data, "#")
	current, err := strconv.Atoi(args[1])
	if err != nil {
		log.Errorf("Parse current number error: %v", err)
		return
	}

	ListRecord(current, args[0], callback.Message.Chat.ID, callback.From.ID, bot)
}

// ListRecord current >= 1
func ListRecord(current int, modelType string, chatID int64, userID int64, bot *BasicTGBot) {
	index := current - 1
	var chattable tgbotapi.Chattable
	switch modelType {
	case string(model.Photo):
		{
			list := cache.GetTypeList(model.Photo, userID)
			if len(list) <= index || index < 0 {
				bot.SendChatMessage("没有记录的数据", chatID)
				return
			}

			files := make([]interface{}, 0)

			var batchForwardDate int64 = 0
			for i := index; i < len(list); i++ {
				fileModel := list[i]
				// 合并同时间 Photo
				if batchForwardDate == 0 {
					batchForwardDate = fileModel.ForwardDate
				} else if batchForwardDate != fileModel.ForwardDate {
					break
				}
				photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(fileModel.FileID))
				photo.Caption = fileModel.Caption
				files = append(files, photo)
				current++
			}

			chattable = tgbotapi.NewMediaGroup(chatID, files)
			bot.Send(chattable)

			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("当前套图保存时间: %s", time.Unix(batchForwardDate, 0).Format("2006-01-02 15:04:05")))
			// 上面 for 循环每一遍都加了偏移, 整体偏移多出1, 此处-1
			msg.ReplyMarkup = generatePhotosInlineKeyboard(string(model.Photo), current-1, len(list))
			bot.Send(msg)
			break
		}
	default:
		{
			log.Warnf("Unknown callback data '%s', is it up to date ?", modelType)
			return
		}
	}
}
