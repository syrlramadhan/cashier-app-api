package repositories

import (
	"github.com/syrlramadhan/cashier-app/models"
	"gorm.io/gorm"
)

type SettingRepository interface {
	FindAll() ([]models.Setting, error)
	FindByKey(key string) (*models.Setting, error)
	Create(setting *models.Setting) error
	Update(setting *models.Setting) error
	UpdateByKey(key string, value string) error
	Delete(id uint) error
}

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) FindAll() ([]models.Setting, error) {
	var settings []models.Setting
	err := r.db.Find(&settings).Error
	return settings, err
}

func (r *settingRepository) FindByKey(key string) (*models.Setting, error) {
	var setting models.Setting
	err := r.db.Where("`key` = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *settingRepository) Create(setting *models.Setting) error {
	return r.db.Create(setting).Error
}

func (r *settingRepository) Update(setting *models.Setting) error {
	return r.db.Save(setting).Error
}

func (r *settingRepository) UpdateByKey(key string, value string) error {
	return r.db.Model(&models.Setting{}).Where("`key` = ?", key).Update("value", value).Error
}

func (r *settingRepository) Delete(id uint) error {
	return r.db.Delete(&models.Setting{}, id).Error
}
