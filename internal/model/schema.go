package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// gorm dto

type FileType string

const (
	// Text 纯文本
	Text     FileType = "Text"
	Photo    FileType = "Photo"
	Voice    FileType = "Voice"
	Video    FileType = "Video"
	Document FileType = "Document"
)

func (ft FileType) Value() (driver.Value, error) {
	return string(ft), nil
}

// Scan 实现 sql.Scanner 接口，从数据库读取值并将其转换回 FileType
func (ft *FileType) Scan(value interface{}) error {
	if value == nil {
		*ft = ""
		return nil
	}

	dbValue, ok := value.(string)
	if !ok {
		return errors.New("invalid file type value from database")
	}

	*ft = FileType(dbValue)
	return nil
}

type FileModel struct {
	gorm.Model
	// 转发消息的 ID
	MessageID int64 `gorm:"index"`
	// 原消息发送时间 (同一 Date 的消息被认为是同一批消息)
	ForwardDate int64 `gorm:"not null"`
	// 转发者 UserId
	UserID int64 `gorm:"index"`
	// 转发者 UserName
	UserName string `gorm:"default:''"`
	// 同一批次消息中由 From 用户发送的文本
	Text string `gorm:"default:''"`
	// 转发消息的原 json 格式数据
	Json datatypes.JSONType[*tgbotapi.Message] `gorm:"not null"`

	// 文件类型 nullable
	FileType FileType `gorm:"default:''"`
	// 文件本地储存路径 nullable
	FilePath string `gorm:"default:''"`
	// 文件大小
	FileSize int64 `gorm:"default:0"`
	// 文件ID
	FileID string `gorm:"default:''"`
}

func (m *FileModel) IsValid() bool {
	switch m.FileType {
	case "":
		{
			return false
		}
	case Text:
		{
			return m.Text != ""
		}
	case Photo:
	case Voice:
	case Video:
	case Document:
		{
			return m.FilePath != "" && m.FileSize != 0
		}
	}
	return true
}

// GetForwardFrom 获取原消息发送者的信息
func (m *FileModel) GetForwardFrom() string {
	message := m.Json.Data()
	if message.ForwardFromChat != nil {
		chat := message.ForwardFromChat
		return fmt.Sprintf("%s(%s)", chat.Title, chat.Type)
	} else if message.ForwardFrom != nil {
		user := message.ForwardFrom
		return fmt.Sprintf("%s %s@%s", user.FirstName, user.LastName, user.UserName)
	} else if message.ForwardSenderName != "" {
		return message.ForwardSenderName
	}
	return "[Unknown]"
}
