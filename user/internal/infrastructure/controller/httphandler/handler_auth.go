package httphandler

import (
	"chatterbox/user/internal/application/usecase"
	"chatterbox/user/internal/domain/valueobject"
	"chatterbox/user/internal/infrastructure/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg *config.Config

	registerUC     *usecase.RegisterUseCase
	loginUC        *usecase.LoginUseCase
	logoutUC       *usecase.LogoutUseCase
	refreshTokenUC *usecase.RefreshTokenUseCase
}

func NewAuthHandler(
	cfg *config.Config,
	registerUC *usecase.RegisterUseCase,
	loginUC *usecase.LoginUseCase,
	logoutUC *usecase.LogoutUseCase,
	refreshTokenUC *usecase.RefreshTokenUseCase,
) *AuthHandler {
	return &AuthHandler{
		cfg:            cfg,
		registerUC:     registerUC,
		loginUC:        loginUC,
		logoutUC:       logoutUC,
		refreshTokenUC: refreshTokenUC,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	res, err := h.registerUC.Execute(
		c.Request.Context(),
		usecase.RegisterCommand{
			Email:       email,
			Username:    req.Username,
			DisplayName: req.DisplayName,
			Password:    req.Password,
		},
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, RegisterResponse{ID: res.UserID.String()})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	res, err := h.loginUC.Execute(
		c.Request.Context(),
		usecase.LoginCommand{
			Email:    email,
			Password: req.Password,
		},
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	c.SetCookie(
		"access_token",
		res.AccessToken,
		int(h.cfg.JWT.AccessTTL.Seconds()),
		"/",
		h.cfg.HttpServer.HostDomain,
		false,
		true,
	)
	c.SetCookie(
		"refresh_token",
		res.RefreshToken,
		int(h.cfg.JWT.RefreshTTL.Seconds()),
		"/",
		h.cfg.HttpServer.HostDomain,
		false,
		true,
	)

	c.JSON(
		http.StatusOK,
		AuthResponse{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
		},
	)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "refresh token not found"})
		return
	}

	err = h.logoutUC.Execute(c.Request.Context(), usecase.LogoutCommand{RefreshToken: refreshToken})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "refresh token not found"})
		return
	}

	res, err := h.refreshTokenUC.Execute(
		c.Request.Context(),
		usecase.RefreshTokenCommand{
			RefreshToken: refreshToken,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.SetCookie(
		"access_token",
		res.AccessToken,
		int(h.cfg.JWT.AccessTTL.Seconds()),
		"/",
		h.cfg.HttpServer.HostDomain,
		false,
		true,
	)
	c.SetCookie(
		"refresh_token",
		res.RefreshToken,
		int(h.cfg.JWT.RefreshTTL.Seconds()),
		"/",
		h.cfg.HttpServer.HostDomain,
		false,
		true,
	)

	c.JSON(
		http.StatusOK,
		AuthResponse{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
		},
	)
}
