package db

import (
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
	log "github.com/sirupsen/logrus"
)

func InsertFile(fileModel *model.FileModel) error {
	result := Instance.Create(fileModel)
	if result.Error != nil {
		return result.Error
	}
	log.Debugf("RowsAffected: %d", result.RowsAffected)
	return nil
}

func QueryAllFile() ([]*model.FileModel, error) {
	fileModels := make([]*model.FileModel, 0)
	result := Instance.Find(&fileModels)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Debugf("Query %d records", result.RowsAffected)
	return fileModels, nil
}

func InsertUser(userId int64) error {
	result := Instance.Create(&model.UserModel{UserID: userId})
	if result.Error != nil {
		return result.Error
	}
	log.Debugf("RowsAffected: %d", result.RowsAffected)
	return nil
}

func ExistUser(userID int64) (bool, error) {
	var user *model.UserModel
	result := Instance.Where("user_id = ?", userID).Find(&user)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
