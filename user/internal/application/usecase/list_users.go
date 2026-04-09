package usecase

import (
	"chatterbox/user/internal/application/dto"
	"chatterbox/user/internal/domain/port"
	"chatterbox/user/internal/domain/valueobject"
	"context"
)

type ListUsersCommand struct {
	UserIDs []valueobject.UserID
}

type ListUsersResult struct {
	Users []*dto.User
}

type ListUsersUseCase struct {
	users port.UserRepo
}

func NewListUsersUseCase(
	users port.UserRepo,
) *ListUsersUseCase {
	return &ListUsersUseCase{
		users: users,
	}
}

func (uc *ListUsersUseCase) Execute(
	ctx context.Context,
	cmd ListUsersCommand,
) (*ListUsersResult, error) {
	users, err := uc.users.FindManyByUserIDs(ctx, cmd.UserIDs)
	if err != nil {
		return nil, err
	}

	usersRes := make([]*dto.User, len(users))
	for i, user := range users {
		usersRes[i] = &dto.User{
			ID:          user.ID().String(),
			Email:       user.Email().String(),
			Username:    user.Username(),
			DisplayName: user.DisplayName(),
			Status:      user.Status().Int(),
		}
	}

	return &ListUsersResult{
		Users: usersRes,
	}, nil
}
