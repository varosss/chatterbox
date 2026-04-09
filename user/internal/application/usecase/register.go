package usecase

import (
	"chatterbox/user/internal/domain/entity"
	"chatterbox/user/internal/domain/port"
	"chatterbox/user/internal/domain/valueobject"
	"context"
	"errors"
)

type RegisterCommand struct {
	Email       valueobject.Email
	Username    string
	DisplayName string
	Password    string
}

type RegisterResult struct {
	UserID valueobject.UserID
}

type RegisterUseCase struct {
	users     port.UserRepo
	passwords port.PasswordHasher
}

func NewRegisterUseCase(
	users port.UserRepo,
	passwords port.PasswordHasher,
) *RegisterUseCase {
	return &RegisterUseCase{
		users:     users,
		passwords: passwords,
	}
}

func (uc *RegisterUseCase) Execute(
	ctx context.Context,
	cmd RegisterCommand,
) (*RegisterResult, error) {
	if uc.users.ExistsByEmail(ctx, cmd.Email) {
		return nil, errors.New("user already exists")
	}

	hash, err := uc.passwords.Hash(cmd.Password)
	if err != nil {
		return nil, err
	}

	passwordHash, err := valueobject.NewPasswordHash(hash)
	if err != nil {
		return nil, err
	}

	u := entity.NewUser(
		cmd.Email,
		cmd.Username,
		cmd.DisplayName,
		passwordHash,
	)

	if err := uc.users.Save(ctx, u); err != nil {
		return nil, err
	}

	return &RegisterResult{
		UserID: u.ID(),
	}, nil
}
