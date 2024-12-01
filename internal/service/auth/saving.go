package auth

import (
	"account-service/internal/domain/users"
	"account-service/pkg/log"
	"context"
	"database/sql"
	"errors"
	"go.uber.org/zap"
)

func (s *Service) SaveUserSavings(ctx context.Context, saving users.Savings) (err error) {
	logger := log.LoggerFromContext(ctx).Named("SaveUserSavings")

	data, err := s.userRepository.GetUserSavings(ctx, saving.UserID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		_, err = s.userRepository.CreateUserSavings(ctx, saving)
		if err != nil {
			logger.Error("failed to save user savings", zap.Error(err))
			return
		}
	case errors.Is(err, nil):
		saving.ID = data.ID
		err = s.userRepository.UpdateUserSavings(ctx, saving)
		if err != nil {
			logger.Error("failed to save user savings", zap.Error(err))
			return
		}
	default:
		logger.Error("failed to get user savings", zap.Error(err))
		return
	}

	return
}

func (s *Service) GetUserSaving(ctx context.Context, userID string) (dest users.Savings, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetUserSaving")

	dest, err = s.userRepository.GetUserSavings(ctx, userID)
	if err != nil {
		logger.Error("failed to get user savings", zap.Error(err))
		return
	}

	return
}
