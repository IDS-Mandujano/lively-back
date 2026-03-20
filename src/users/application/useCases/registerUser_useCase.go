package usecases

import (
	"context"
	"errors"
	"lively-backend/src/users/domain/entity"
	"lively-backend/src/users/domain/repository"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserUseCase struct {
	userRepo repository.IUserRepository
}

func NewRegisterUserUseCase(repo repository.IUserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo: repo,
	}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*entity.User, error) {

	if input.Username == "" || input.Email == "" || input.Password == "" {
		return nil, errors.New("todos los campos son obligatorios")
	}

	if len(input.Password) < 6 {
		return nil, errors.New("la contraseña debe tener al menos 6 caracteres")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error al encriptar la contraseña")
	}

	user := &entity.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
	}

	err = uc.userRepo.Save(ctx, user)
	if err != nil {
		return nil, errors.New("el correo o nombre de usuario ya están registrados")
	}

	return user, nil
}
