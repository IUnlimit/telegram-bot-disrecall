package db

import (
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
	log "github.com/sirupsen/logrus"
)

func Insert(fileModel *model.FileModel) error {
	result := Instance.Create(fileModel)
	if result.Error != nil {
		return result.Error
	}
	log.Debugf("RowsAffected: %d", result.RowsAffected)
	return nil
}
