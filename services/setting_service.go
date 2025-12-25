package services

import (
	"errors"

	"github.com/syrlramadhan/cashier-app/dto"
	"github.com/syrlramadhan/cashier-app/models"
	"github.com/syrlramadhan/cashier-app/repositories"
)

type SettingService struct {
	settingRepo repositories.SettingRepository
}

func NewSettingService(settingRepo repositories.SettingRepository) *SettingService {
	return &SettingService{settingRepo: settingRepo}
}

func (s *SettingService) GetAllSettings() ([]dto.SettingResponse, error) {
	settings, err := s.settingRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []dto.SettingResponse
	for _, setting := range settings {
		response = append(response, dto.SettingResponse{
			Key:   setting.Key,
			Value: setting.Value,
		})
	}

	return response, nil
}

func (s *SettingService) GetSettingByKey(key string) (*dto.SettingResponse, error) {
	setting, err := s.settingRepo.FindByKey(key)
	if err != nil {
		return nil, errors.New("setting not found")
	}

	response := &dto.SettingResponse{
		Key:   setting.Key,
		Value: setting.Value,
	}

	return response, nil
}

func (s *SettingService) UpdateSetting(req *dto.UpdateSettingRequest) (*dto.SettingResponse, error) {
	_, err := s.settingRepo.FindByKey(req.Key)
	if err != nil {
		// Create if not exists
		setting := &models.Setting{
			Key:   req.Key,
			Value: req.Value,
		}
		err = s.settingRepo.Create(setting)
		if err != nil {
			return nil, errors.New("failed to create setting")
		}
	} else {
		err = s.settingRepo.UpdateByKey(req.Key, req.Value)
		if err != nil {
			return nil, errors.New("failed to update setting")
		}
	}

	return &dto.SettingResponse{
		Key:   req.Key,
		Value: req.Value,
	}, nil
}
