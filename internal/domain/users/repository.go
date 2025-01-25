package users

import "context"

type Repository interface {
	CreateUser(ctx context.Context, data User) (string, error)
	CreateDependency(ctx context.Context, data Dependency) (string, error)
	CreateUserDependency(ctx context.Context, data UserDependency) (string, error)
	CreateSobrietyTracking(ctx context.Context, data SobrietyTracking) (string, error)

	GetUsers(ctx context.Context) ([]User, error)

	GetUserDependency(ctx context.Context, id, dId string) (dest UserDependency, err error)
	GetUserDependencies(ctx context.Context, id string) (dest []UserDependency, err error)
	GetDependencies(ctx context.Context) ([]Dependency, error)
	GetUserByEmailOrLogin(ctx context.Context, email string, login string) (User, error)
	GetDependency(ctx context.Context, id string) (Dependency, error)
	GetSobrietyTracking(ctx context.Context, userID string) (SobrietyTracking, error)
	GetUserByAny(ctx context.Context, login string) (dest User, err error)

	GetUserSavings(ctx context.Context, userId string) (dest Savings, err error)
	UpdateUserSavings(ctx context.Context, data Savings) (err error)
	CreateUserSavings(ctx context.Context, data Savings) (id string, err error)
}
