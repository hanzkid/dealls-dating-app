package handler

import (
	"main/helpers"
	"main/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ProfileRequest struct {
	Description *string `json:"description,omitempty"`
	Picture     *string `json:"picture,omitempty"`
}

type UserHandler struct {
	userRepository    repository.UserRepositoryInterface
	profileRepository repository.ProfileRepositoryInterface
}

func NewUserHandler(userRepository repository.UserRepositoryInterface, profileRepository repository.ProfileRepositoryInterface) *UserHandler {
	return &UserHandler{
		userRepository,
		profileRepository,
	}
}

func (h *UserHandler) Me(c echo.Context) error {
	userId := c.Get("user_id").(int)
	user, err := h.userRepository.FindByID(userId)
	if err != nil {
		helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal server error")
		return nil
	}

	helpers.ResponseWithSuccess(c, http.StatusOK, user)
	return nil
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userId := c.Get("user_id").(int)

	var profileRequest ProfileRequest
	if err := c.Bind(&profileRequest); err != nil {
		helpers.ResponseWithError(c, http.StatusBadRequest, "Invalid request")
		return nil
	}

	profile, err := h.profileRepository.FindByUserID(userId)
	if err != nil {
		helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal server error")
		return nil
	}
	if profile == nil {
		helpers.ResponseWithError(c, http.StatusNotFound, "Profile not found")
		return nil
	}

	if profileRequest.Description == nil && profileRequest.Picture == nil {
		helpers.ResponseWithError(c, http.StatusBadRequest, "Nothing to update")
		return nil
	}

	if profileRequest.Description != nil {
		profile.Description = *profileRequest.Description
	}
	if profileRequest.Picture != nil {
		profile.Picture = *profileRequest.Picture
	}
	_, err = h.profileRepository.Save(profile)
	if err != nil {
		helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal server error")
		return nil
	}

	helpers.ResponseWithSuccess(c, http.StatusOK, map[string]interface{}{"message": "Profile updated"})
	return nil
}

func (h *UserHandler) PurchasePremium(c echo.Context) error {

	checkActiveSubscription, err := h.userRepository.CheckSubscription(c)
	if err != nil {
		helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal server error")
		return nil
	}

	if checkActiveSubscription {
		helpers.ResponseWithError(c, http.StatusBadRequest, "You already have an active subscription")
		return nil
	}

	user, err := h.userRepository.Subscribe(c)
	if err != nil {
		helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal server error")
		return nil
	}

	helpers.ResponseWithSuccess(c, http.StatusOK, map[string]interface{}{"message": "Successfully purchased premium valid until " + user.Subscription.ValidUntil.String()})
	return nil
}
