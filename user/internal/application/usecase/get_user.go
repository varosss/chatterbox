package usecase

import (
	"chatterbox/user/internal/application/dto"
	"chatterbox/user/internal/domain/port"
	"chatterbox/user/internal/domain/valueobject"
	"context"
)

type GetUserCommand struct {
	UserID valueobject.UserID
}

type GetUserResult struct {
	User *dto.User
}

type GetUserUseCase struct {
	users port.UserRepo
}

func NewGetUserUseCase(
	users port.UserRepo,
) *GetUserUseCase {
	return &GetUserUseCase{
		users: users,
	}
}

func (uc *GetUserUseCase) Execute(
	ctx context.Context,
	cmd GetUserCommand,
) (*GetUserResult, error) {
	user, err := uc.users.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	return &GetUserResult{
		User: &dto.User{
			ID:          user.ID().String(),
			Email:       user.Email().String(),
			Username:    user.Username(),
			DisplayName: user.DisplayName(),
			Status:      user.Status().Int(),
		},
	}, nil
}
