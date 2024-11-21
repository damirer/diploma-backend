package users

import (
	"account-service/internal/domain/grant"
	"time"
)

type Dependency struct {
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	Name string `db:"name" json:"name"`
	ID   string `db:"id,primarykey" json:"id"`
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
	ID           string `db:"id,primarykey"`
	UserID       string `db:"user_id"`
	DependencyID string `db:"dependency_id" json:"id"`
}

type SobrietyTracking struct {
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
	ID             string    `db:"id,primarykey" json:"id" json:"id,omitempty"`
	UserID         string    `db:"user_id"`
	StartDate      time.Time `db:"start_date" json:"startDate"`
	DaysPerWeek    int       `db:"days_per_week" json:"daysPerWeek,omitempty"`
	DrinksPerDay   int       `db:"drinks_per_day" json:"drinksPerDay,omitempty"`
	MotivationText string    `db:"motivation_text" json:"motivationText,omitempty"`
}

func ParseFromAuth(req grant.Request) User {
	return User{
		Email:    req.Login,
		Password: req.Password,
	}
}
