package postgres

import (
	"context"
	"fmt"
	"strings"
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

func (r *Repository) GetUserByEmailOrLogin(ctx context.Context, email string, login string) (dest users.User, err error) {
	query := `
			SELECT created_at, updated_at, id, email, name, password
			FROM users
			WHERE email = $1 OR name = $2;`

	args := []interface{}{email, login}

	err = r.db.GetContext(ctx, &dest, query, args...)
	return
}

func (r *Repository) GetUserByAny(ctx context.Context, login string) (dest users.User, err error) {
	query := `
			SELECT created_at, updated_at, id, email, name, password
			FROM users
			WHERE email = $1 OR name = $1;`

	args := []interface{}{login}

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

func (r *Repository) CreateUserSavings(ctx context.Context, data users.Savings) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
			INSERT INTO user_savings (user_id,	money,	start_date)
			VALUES ($1, $2,$3)
			RETURNING id;`

	args := []interface{}{data.UserID, *data.Money, *data.StartDate}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	return
}

func (r *Repository) GetUserSavings(ctx context.Context, userId string) (dest users.Savings, err error) {
	query := `
			SELECT id, user_id,	money,	start_date
			FROM user_savings
			WHERE user_id = $1;`

	args := []interface{}{userId}

	err = r.db.GetContext(ctx, &dest, query, args...)

	return
}

func (r *Repository) UpdateUserSavings(ctx context.Context, data users.Savings) (err error) {
	sets, args := r.prepareArgs(data)
	if len(args) > 0 {
		args = append(args, data.ID)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE user_savings SET %s WHERE id=$%d", strings.Join(sets, ", "), len(args))

		_, err = r.db.ExecContext(ctx, query, args...)
	}

	return
}

func (r *Repository) prepareArgs(data users.Savings) (sets []string, args []any) {
	if data.StartDate != nil {
		args = append(args, data.StartDate)
		sets = append(sets, fmt.Sprintf("start_date=$%d", len(args)))
	}
	if data.Money != nil {
		args = append(args, data.Money)
		sets = append(sets, fmt.Sprintf("money=$%d", len(args)))
	}

	return
}
