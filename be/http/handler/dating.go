package handler

import (
	"net/http"

	"main/entity"
	"main/helpers"
	"main/repository"

	"github.com/labstack/echo/v4"
)

type DatingHandler struct {
	profileRepository repository.ProfileRepositoryInterface
	matchRepository   repository.MatchRepositoryInterface
}

type SwipeRequest struct {
	ProfileID int  `json:"profile_id"`
	Swipe     bool `json:"swipe"`
}

func NewDatingHandler(profileRepository repository.ProfileRepositoryInterface, matchRepository repository.MatchRepositoryInterface) *DatingHandler {
	return &DatingHandler{
		profileRepository: profileRepository,
		matchRepository:   matchRepository,
	}
}

func (h *DatingHandler) Profile(c echo.Context) error {
	checkDailyLimit, err := h.matchRepository.CheckDailyLimit(c)
	if err != nil {
		helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
		return nil
	}
	if !checkDailyLimit {
		helpers.ResponseWithError(c, http.StatusForbidden, "Daily limit reached")
		return nil
	}
	profile, err := h.profileRepository.GetRandomProfile(c)
	if err != nil {
		helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
		return nil
	}

	err = h.profileRepository.SaveViewLog(c, int(profile.ID))
	if err != nil {
		helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
		return nil
	}
	helpers.ResponseWithSuccess(c, http.StatusOK, profile)
	return nil
}

func (h *DatingHandler) SwipedProfile(c echo.Context) error {
	var req SwipeRequest
	if err := c.Bind(&req); err != nil {
		helpers.ResponseWithError(c, http.StatusBadRequest, "Invalid request")
		return nil
	}
	partnerId := req.ProfileID
	profileId := c.Get("profile_id").(int)

	_, err := h.profileRepository.FindByID(profileId)
	if err != nil {
		helpers.ResponseWithError(c, http.StatusNotFound, "Profile not found")
		return nil
	}

	// if user swipe left, reject the match for profile that swiped right
	if !req.Swipe {
		pendingMatch, err := h.matchRepository.CheckPendingMatch(partnerId, profileId)
		if err != nil {
			helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
			return nil
		}
		if pendingMatch == nil {
			err := h.Profile(c)
			if err != nil {
				c.Logger().Error(err)
				helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
			}
			return nil
		}
		err = h.matchRepository.RejectMatch(profileId, partnerId)
		if err != nil {
			c.Logger().Error(err)
			helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
			return nil
		}

		err = h.Profile(c)
		if err != nil {
			c.Logger().Error(err)
			helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
			return nil
		}
		return nil
	}

	// check if user already swiped
	userSwiped, err := h.matchRepository.CheckMatch(profileId, partnerId)
	if err != nil {
		helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
		return nil
	}
	if userSwiped != nil {
		if userSwiped.Status == entity.StatusPending {
			helpers.ResponseWithError(c, http.StatusConflict, "Already swiped")
			return nil
		}
		helpers.ResponseWithError(c, http.StatusConflict, "Already matched")
		return nil
	}

	// check if user swiped right and the profile that swiped right also swiped right
	partnerSwiped, err := h.matchRepository.CheckMatch(partnerId, profileId)
	if err != nil {
		helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
		return nil
	}

	// if partner already swiped right, accept the match
	if partnerSwiped != nil {
		err := h.matchRepository.AcceptMatch(profileId, partnerId)
		if err != nil {
			helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
			return nil
		}
	} else {
		err := h.matchRepository.CreateMatch(profileId, partnerId)
		if err != nil {
			helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
			return nil
		}
	}
	err = h.Profile(c)
	if err != nil {
		c.Logger().Error(err)
		helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
	}
	return nil
}

func (h *DatingHandler) MatchList(c echo.Context) error {
	profileId := c.Get("profile_id").(int)
	matches, err := h.matchRepository.FindMatchByProfileID(profileId)
	if err != nil {
		helpers.ResponseWithError(c, http.StatusInternalServerError, "Internal Server Error")
		return nil
	}
	helpers.ResponseWithSuccess(c, http.StatusOK, matches)
	return nil
}
