package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IUnlimit/telegram-bot-disrecall/internal/cache"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

// InputChanMap chatID: chan
var InputChanMap map[int64]chan string

var supportTypes map[string]func(tgbotapi.RequestFileData, *model.FileModel) interface{}

// OnCommandCallback 共计三种回调数据类型
// 1. FileType#index 跳转页面
// 2. /command 便捷
// 3. >FileType#1-total 捕获指定范围内输入数字, 随后跳转页面
func OnCommandCallback(callback *tgbotapi.CallbackQuery, bot *BasicTGBot) {
	data := callback.Data
	log.Debugf("Callback query received from user %d with data %s", callback.From.ID, data)

	if strings.HasPrefix(data, "/") {
		// 回调指令
		callback.Message.Text = data
		OnCommand(callback.Message, bot)
		return
	} else if strings.HasPrefix(data, ">") {
		// 回调捕获
		onSearch(callback, bot)
		return
	}

	args := strings.Split(data, "#")
	current, err := strconv.Atoi(args[1])
	if err != nil {
		log.Errorf("Parse current number error: %v", err)
		return
	}

	ListRecord(current, args[0], callback.Message.Chat.ID, callback.From.ID, bot)
}

// ListRecord current >= 1
func ListRecord(current int, modelTypeStr string, chatID int64, userID int64, bot *BasicTGBot) {
	if _, ok := supportTypes[modelTypeStr]; !ok {
		bot.SendChatMessage(fmt.Sprintf("不支持查询的数据类型: %s", modelTypeStr), chatID)
		return
	}
	generateFunc := supportTypes[modelTypeStr]

	index := current - 1
	var modelType model.FileType
	_ = modelType.Scan(modelTypeStr)

	list := cache.GetTypeList(modelType, userID)
	if len(list) <= index || index < 0 {
		bot.SendChatMessage("没有记录的数据", chatID)
		return
	}

	files := make([]interface{}, 0)

	var batchForwardMediaGroupID = ""
	for i := index; i < len(list); i++ {
		fileModel := list[i]
		// 合并同时间 Media
		if batchForwardMediaGroupID == "" {
			batchForwardMediaGroupID = fileModel.MediaGroupID
		} else if batchForwardMediaGroupID != fileModel.MediaGroupID {
			break
		}
		// return &
		media := generateFunc(tgbotapi.FileID(fileModel.FileID), fileModel)
		files = append(files, media)
		current++
	}

	var chattable tgbotapi.Chattable
	if len(files) == 1 && isTextMessage(files[0]) {
		config := files[0].(*tgbotapi.MessageConfig)
		config.ChatID = chatID
		chattable = config
	} else {
		chattable = tgbotapi.NewMediaGroup(chatID, files)
	}
	bot.Send(chattable)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("当前 [%s] 保存id: %s\n\t——%s",
		modelTypeStr,
		batchForwardMediaGroupID,
		list[index].GetForwardFrom(),
	))
	// 上面 for 循环每一遍都加了偏移, 整体偏移多出1, 此处-1
	msg.ReplyMarkup = generateQueryInlineKeyboard(modelTypeStr, current-1, len(list))
	bot.Send(msg)
}

func onSearch(callback *tgbotapi.CallbackQuery, bot *BasicTGBot) {
	data := callback.Data
	chatID := callback.Message.Chat.ID

	args := strings.Split(data[1:], "#")
	start, end := strRange(args[1])
	channel := make(chan string)
	InputChanMap[chatID] = channel
	bot.SendChatMessage(fmt.Sprintf("请输入你要跳转到的页面 (%d-%d)", start, end), chatID)

	for {
		inputStr := <-channel
		jump, err := strconv.Atoi(inputStr)
		if err != nil || (jump < start || jump > end) {
			bot.SendChatMessage("请输入正确的页数", chatID)
			continue
		}
		delete(InputChanMap, chatID)
		ListRecord(jump, args[0], chatID, callback.From.ID, bot)
		break
	}
}

func isTextMessage(entity interface{}) bool {
	if _, ok := entity.(*tgbotapi.MessageConfig); ok {
		return true
	}
	return false
}

// str '%d-%d'
func strRange(str string) (int, int) {
	splits := strings.Split(str, "-")
	start, _ := strconv.Atoi(splits[0])
	end, _ := strconv.Atoi(splits[1])
	return start, end
}
