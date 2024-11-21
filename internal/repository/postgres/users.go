package postgres

import (
	"context"
	"time"

	"account-service/internal/domain/users"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, data users.User) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
			INSERT INTO users (email, name, password)
			VALUES ($1, $2,$3)
			RETURNING id;`

	args := []interface{}{data.Email, data.Name, data.Password}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	return
}

func (r *Repository) CreateDependency(ctx context.Context, data users.Dependency) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
			INSERT INTO dependencies (name)
			VALUES ($1)
			RETURNING id;`

	args := []interface{}{data.Name}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	return
}

func (r *Repository) CreateUserDependency(ctx context.Context, data users.UserDependency) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
			INSERT INTO user_dependencies (user_id, dependency_id)
			VALUES ($1, $2)
			RETURNING id;`

	args := []interface{}{data.UserID, data.DependencyID}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	return
}

func (r *Repository) GetUserDependency(ctx context.Context, id, dId string) (dest users.UserDependency, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `SELECT   user_id, dependency_id
			FROM user_dependencies
	WHERE user_id=$1 AND dependency_id=$2`

	args := []interface{}{id, dId}

	err = r.db.GetContext(ctx, &dest, query, args...)
	return
}

func (r *Repository) CreateSobrietyTracking(ctx context.Context, data users.SobrietyTracking) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
			INSERT INTO sobriety_tracking (user_id, start_date, days_per_week, drinks_per_day, motivation_text)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id;`

	args := []interface{}{data.UserID, data.StartDate, data.DaysPerWeek, data.DrinksPerDay, data.MotivationText}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	return
}

func (r *Repository) GetSobrietyTracking(ctx context.Context, userID string) (dest users.SobrietyTracking, err error) {
	query := `SELECT id, user_id, start_date, days_per_week, drinks_per_day, motivation_text, FROM sobriety_tracking WHERE user_id = $1 ORDER BY start_date DESC LIMIT 1`

	if err = r.db.GetContext(ctx, &dest, query, userID); err != nil {
		return
	}

	return
}

func (r *Repository) GetDependencies(ctx context.Context) (dest []users.Dependency, err error) {
	query := `
			SELECT   created_at, updated_at,id,name
			FROM dependencies
			ORDER BY created_at DESC;`

	err = r.db.SelectContext(ctx, &dest, query)
	return
}

func (r *Repository) GetUserByEmail(ctx context.Context, id string) (dest users.User, err error) {
	query := `
			SELECT created_at, updated_at, id, email, name, password
			FROM users
			WHERE email = $1;`

	args := []interface{}{id}

	err = r.db.GetContext(ctx, &dest, query, args...)
	return
}

func (r *Repository) GetDependency(ctx context.Context, id string) (dest users.Dependency, err error) {
	query := `
			SELECT id, name, created_at, updated_at
			FROM dependencies
			WHERE id = $1;`

	args := []interface{}{id}

	err = r.db.GetContext(ctx, &dest, query, args...)
	return
}
