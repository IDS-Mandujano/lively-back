package mysql

import (
	"context"
	"database/sql"
	"errors"
	"lively-backend/src/users/domain/entity"
)

type MySQLUserRepository struct {
	db *sql.DB // Recibimos la conexión global a la base de datos
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{
		db: db,
	}
}

// Save inserta un nuevo usuario en la base de datos
func (repo *MySQLUserRepository) Save(ctx context.Context, user *entity.User) error {
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"

	// Ejecutamos el query protegiéndonos de inyecciones SQL usando los signos de interrogación
	result, err := repo.db.ExecContext(ctx, query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return err // Puede fallar si el correo o el username ya existen (son UNIQUE)
	}

	// Obtenemos el ID que MySQL generó automáticamente y se lo asignamos a la entidad
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)

	return nil
}

// FindByEmail busca a un usuario en la DB usando su correo electrónico
func (repo *MySQLUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := "SELECT id, username, email, password_hash, created_at FROM users WHERE email = ?"
	row := repo.db.QueryRowContext(ctx, query, email)

	var user entity.User
	// Escaneamos el resultado de MySQL hacia las variables de Go
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("credenciales incorrectas o usuario no encontrado")
		}
		return nil, err
	}

	return &user, nil
}
