package mysql

import (
	"context"
	"database/sql"
	"errors"
	"lively-backend/src/users/domain/entity"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{
		db: db,
	}
}

func (repo *MySQLUserRepository) Save(ctx context.Context, user *entity.User) error {
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"

	result, err := repo.db.ExecContext(ctx, query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)

	return nil
}

func (repo *MySQLUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := "SELECT id, username, email, password_hash, created_at FROM users WHERE email = ?"
	row := repo.db.QueryRowContext(ctx, query, email)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("credenciales incorrectas o usuario no encontrado")
		}
		return nil, err
	}

	return &user, nil
}
