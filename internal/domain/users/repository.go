package users

import "context"

type Repository interface {
	CreateUser(ctx context.Context, data User) (string, error)
	CreateDependency(ctx context.Context, data Dependency) (string, error)
	CreateUserDependency(ctx context.Context, data UserDependency) (string, error)
	CreateSobrietyTracking(ctx context.Context, data SobrietyTracking) (string, error)

	GetUserByEmailAndPassword(ctx context.Context, email, password string) (dest User, err error)
	GetUserByEmail(ctx context.Context, id string) (User, error)
	GetDependency(ctx context.Context, id string) (Dependency, error)
	GetSobrietyTracking(ctx context.Context, userID string) (SobrietyTracking, error)
}