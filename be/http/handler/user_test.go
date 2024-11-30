package handler

import (
	"main/entity"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

// FindByEmail implements repository.UserRepositoryInterface.
func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(*entity.User), args.Error(1)
}

// Save implements repository.UserRepositoryInterface.
func (m *MockUserRepository) Save(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	return args.Get(0).(*entity.User), args.Error(1)
}

// Update implements repository.UserRepositoryInterface.
func (m *MockUserRepository) Update(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id int) (*entity.User, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) CheckSubscription(c echo.Context) (bool, error) {
	args := m.Called(c)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Subscribe(c echo.Context) (*entity.User, error) {
	args := m.Called(c)
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestUserHandler_Me(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", 1)

	mockUserRepo := new(MockUserRepository)
	mockProfileRepo := new(MockProfileRepository)
	handler := NewUserHandler(mockUserRepo, mockProfileRepo)

	user := &entity.User{ID: 1, Name: "John Doe"}
	mockUserRepo.On("FindByID", 1).Return(user, nil)

	err := handler.Me(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUserRepo.AssertExpectations(t)
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	e := echo.New()
	payload := `{"description": "New Description", "picture": "new.jpg"}`
	req := httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", 1)

	mockUserRepo := new(MockUserRepository)
	mockProfileRepo := new(MockProfileRepository)
	handler := NewUserHandler(mockUserRepo, mockProfileRepo)

	profile := &entity.Profile{UserID: 1, Description: "Old Description", Picture: "old.jpg"}
	mockProfileRepo.On("FindByUserID", 1).Return(profile, nil)
	mockProfileRepo.On("Save", mock.AnythingOfType("*entity.Profile")).Return(profile, nil)

	err := handler.UpdateProfile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Profile updated")
	mockProfileRepo.AssertExpectations(t)
}

func TestUserHandler_PurchasePremium(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/purchase", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserRepo := new(MockUserRepository)
	mockProfileRepo := new(MockProfileRepository)
	handler := NewUserHandler(mockUserRepo, mockProfileRepo)

	mockUserRepo.On("CheckSubscription", c).Return(false, nil)
	user := &entity.User{ID: 1, Subscription: entity.Subscription{ValidUntil: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)}}
	mockUserRepo.On("Subscribe", c).Return(user, nil)

	err := handler.PurchasePremium(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUserRepo.AssertExpectations(t)
}

func TestUserHandler_PurchasePremiumAlreadyActive(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/purchase", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserRepo := new(MockUserRepository)
	mockProfileRepo := new(MockProfileRepository)
	handler := NewUserHandler(mockUserRepo, mockProfileRepo)

	mockUserRepo.On("CheckSubscription", c).Return(true, nil)

	_ = handler.PurchasePremium(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "You already have an active subscription")
	mockUserRepo.AssertExpectations(t)
}
