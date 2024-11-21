package users

import (
	"account-service/internal/domain/grant"
	"time"
)

type Dependency struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Name     string `db:"name"`
	ID       string `db:"id,primarykey"`
	Password string `db:"password"`
}

type User struct {
	ID        string    `db:"id,primarykey"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserDependency struct {
	ID           string    `db:"id,primarykey"`
	UserID       string    `db:"user_id"`
	DependencyID string    `db:"dependency_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type SobrietyTracking struct {
	ID             string    `db:"id,primarykey"`
	UserID         string    `db:"user_id"`
	StartDate      time.Time `db:"start_date"`
	DaysPerWeek    int       `db:"days_per_week"`
	DrinksPerDay   int       `db:"drinks_per_day"`
	MotivationText string    `db:"motivation_text"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func ParseFromAuth(req grant.Request) User {
	return User{
		Email:    req.Login,
		Password: req.Password,
	}
}
