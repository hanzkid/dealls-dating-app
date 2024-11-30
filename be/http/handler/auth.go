package handler

import (
	"log"
	"main/config"
	"main/entity"
	"main/helpers"
	"main/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	userRepo repository.UserRepositoryInterface
	cfg      *config.Config
}

func NewAuthHandler(userRepo repository.UserRepositoryInterface, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (h *AuthHandler) Login(echoCtx echo.Context) error {
	var req LoginRequest
	if err := echoCtx.Bind(&req); err != nil {
		return echoCtx.JSON(http.StatusBadRequest, err)
	}
	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		helpers.ResponseWithError(echoCtx, http.StatusUnauthorized, "Invalid credentials")
		return nil
	}

	correctPassword := helpers.ComparePassword(user.Password, req.Password)
	if !correctPassword {
		helpers.ResponseWithError(echoCtx, http.StatusUnauthorized, "Invalid credentials")
		return nil
	}

	accessToken, err := helpers.GenerateAccessToken(user, &h.cfg.JWT)

	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		helpers.ResponseWithError(echoCtx, http.StatusInternalServerError, "Invalid credentials")
		return nil
	}

	helpers.ResponseWithSuccess(echoCtx, http.StatusOK, map[string]string{"access_token": *accessToken})
	return nil
}

func (h *AuthHandler) Register(echoCtx echo.Context) error {
	var req RegisterRequest
	if err := echoCtx.Bind(&req); err != nil {
		return echoCtx.JSON(http.StatusBadRequest, err)
	}
	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	createdUser, err := h.userRepo.Save(&user)
	if err != nil {
		helpers.ResponseWithError(echoCtx, http.StatusInternalServerError, "Internal server error")
		return nil
	}
	responseData := RegisterResponse{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
	}
	helpers.ResponseWithSuccess(echoCtx, http.StatusCreated, responseData)
	return nil
}
