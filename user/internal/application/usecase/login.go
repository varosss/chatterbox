package usecase

import (
	"chatterbox/user/internal/domain/entity"
	"chatterbox/user/internal/domain/port"
	"chatterbox/user/internal/domain/valueobject"
	"context"
	"errors"
	"time"
)

type LoginCommand struct {
	Email    valueobject.Email
	Password string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}

type LoginUseCase struct {
	users         port.UserRepo
	passwords     port.PasswordVerifier
	refreshTokens port.RefreshTokenRepo
	signer        port.TokenSigner
	clock         port.Clock
	refreshTTL    time.Duration
}

func NewLoginUseCase(
	users port.UserRepo,
	passwords port.PasswordVerifier,
	refreshTokens port.RefreshTokenRepo,
	signer port.TokenSigner,
	clock port.Clock,
	refreshTTL time.Duration,
) *LoginUseCase {
	return &LoginUseCase{
		users:         users,
		passwords:     passwords,
		refreshTokens: refreshTokens,
		signer:        signer,
		clock:         clock,
		refreshTTL:    refreshTTL,
	}
}

func (uc *LoginUseCase) Execute(
	ctx context.Context,
	cmd LoginCommand,
) (*LoginResult, error) {
	user, err := uc.users.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}

	if !uc.passwords.Compare(user.PasswordHash().String(), cmd.Password) {
		return nil, errors.New("password is invalid")
	}

	now := uc.clock.Now()

	accessJWT, err := uc.signer.SignAccess(user.ID(), now)
	if err != nil {
		return nil, err
	}

	refresh := entity.NewRefreshToken(
		user.ID(),
		now.Add(uc.refreshTTL),
	)

	if err := uc.refreshTokens.Save(ctx, refresh); err != nil {
		return nil, err
	}

	refreshJWT, err := uc.signer.SignRefresh(refresh.ID(), user.ID(), now)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessJWT,
		RefreshToken: refreshJWT,
	}, nil
}
