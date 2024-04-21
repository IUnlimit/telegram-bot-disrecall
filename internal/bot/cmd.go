package bot

import (
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func OnCommand(message *tgbotapi.Message, basic *BasicTGBot) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	switch message.Text {
	case "/start":
		{
			msg.Text = "您好，我是防撤回机器人，您可将需要保存的 文本/图片/语音/视频/文件 转发至此机器人，机器人将会自动将文件存档到本地服务器。若原消息被撤回，则机器人会将存档文件重新上传至该聊天"
			msg.ReplyMarkup = replyKeyboard
			break
		}
	case "/photos":
		{
			ListRecord(1, string(model.Photo), message.Chat.ID, message.From.ID, basic)
			break
		}
	case "/help":
		{
			msg.Text = `帮助菜单
			- /start 初始化聊天
			- /help 获取本菜单`
			break
		}
	default:
		{
			msg.Text = "未知的命令: " + message.Text
		}
	}

	if _, err := basic.API.Send(msg); err != nil {
		log.Error(err)
	}
}

func onPhotos() {

}
