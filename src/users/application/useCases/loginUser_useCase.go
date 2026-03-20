package usecases

import (
	"context"
	"errors"
	"lively-backend/src/users/domain/entity"
	"lively-backend/src/users/domain/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserUseCase struct {
	userRepo repository.IUserRepository
}

func NewLoginUserUseCase(repo repository.IUserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo: repo,
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, input LoginUserInput) (string, *entity.User, error) {
	// 1. Buscamos al usuario por su correo
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return "", nil, errors.New("correo o contraseña incorrectos")
	}

	// 2. Comparamos la contraseña en texto plano con el Hash de la DB
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return "", nil, errors.New("correo o contraseña incorrectos")
	}

	// 3. ¡Contraseña correcta! Fabricamos la "Pulsera VIP" (JWT)
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", nil, errors.New("error interno: llave secreta no configurada")
	}

	// Metemos datos en el token (quién es y cuándo caduca)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // El token dura 3 días
	})

	// Firmamos el token con nuestra llave secreta
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", nil, errors.New("error al generar el token de acceso")
	}

	return tokenString, user, nil
}
