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

func QueryAll() ([]*model.FileModel, error) {
	fileModels := make([]*model.FileModel, 0)
	result := Instance.Find(&fileModels)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Debugf("Query %d records", result.RowsAffected)
	return fileModels, nil
}
