package users

import "context"

type Repository interface {
	CreateUser(ctx context.Context, data User) (string, error)
	CreateDependency(ctx context.Context, data Dependency) (string, error)
	CreateUserDependency(ctx context.Context, data UserDependency) (string, error)
	CreateSobrietyTracking(ctx context.Context, data SobrietyTracking) (string, error)

	GetUserDependency(ctx context.Context, id, dId string) (dest UserDependency, err error)
	GetDependencies(ctx context.Context) ([]Dependency, error)
	GetUserByEmail(ctx context.Context, id string) (User, error)
	GetDependency(ctx context.Context, id string) (Dependency, error)
	GetSobrietyTracking(ctx context.Context, userID string) (SobrietyTracking, error)
}
