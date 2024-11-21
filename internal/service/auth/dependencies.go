package auth

import (
	"account-service/internal/domain/users"
	"account-service/pkg/log"
	"context"
	"database/sql"
	"errors"
	"go.uber.org/zap"
)

func (s *Service) CreateSobrietyTracking(ctx context.Context, req users.SobrietyTracking) (err error) {
	logger := log.LoggerFromContext(ctx).Named("CreateSobrietyTracking")

	_, err = s.userRepository.CreateSobrietyTracking(ctx, req)
	if err != nil {
		logger.Error("failed to create sobriety tracking", zap.Error(err))
	}

	return
}

func (s *Service) GetDependencies(ctx context.Context) (dest []users.Dependency, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetDependencies")

	dest, err = s.userRepository.GetDependencies(ctx)
	if err != nil {
		logger.Error("failed to get dependencies", zap.Error(err))
		return
	}

	return
}

func (s *Service) CreateUserDependencies(ctx context.Context, userId string, deps []users.UserDependency) (err error) {
	logger := log.LoggerFromContext(ctx).Named("CreateUserDependencies")

	for _, dep := range deps {
		_, err = s.userRepository.GetUserDependency(ctx, userId, dep.DependencyID)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			dep.UserID = userId
			dep.ID, err = s.userRepository.CreateUserDependency(ctx, dep)
			if err != nil {
				logger.Error("failed to create dependency", zap.Error(err))
				return
			}
		case errors.Is(err, nil):
			continue
		default:
			logger.Error("failed to get dependencies", zap.Error(err))
			return
		}
	}
	return
}
