package usecase

import (
	"chatterbox/user/internal/domain/port"
	"context"
	"errors"
)

type LogoutCommand struct {
	RefreshToken string
}

type LogoutUseCase struct {
	verifier      port.TokenVerifier
	refreshTokens port.RefreshTokenRepo
}

func NewLogoutUseCase(
	verifier port.TokenVerifier,
	refreshTokens port.RefreshTokenRepo,
) *LogoutUseCase {
	return &LogoutUseCase{
		verifier:      verifier,
		refreshTokens: refreshTokens,
	}
}

func (uc *LogoutUseCase) Execute(
	ctx context.Context,
	cmd LogoutCommand,
) error {
	tokenID, _, err := uc.verifier.VerifyRefresh(cmd.RefreshToken)
	if err != nil {
		return errors.New("invalid token")
	}

	return uc.refreshTokens.Revoke(ctx, tokenID)
}
