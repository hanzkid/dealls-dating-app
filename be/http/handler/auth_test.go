package handler

import (
	"errors"
	"main/config"
	"main/entity"
	"main/helpers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

// FindByID implements repository.UserRepositoryInterface.
func (m *MockUserRepo) FindByID(id int) (*entity.User, error) {
	panic("unimplemented")
}

// Subscribe implements repository.UserRepositoryInterface.
func (m *MockUserRepo) Subscribe(c echo.Context) (*entity.User, error) {
	panic("unimplemented")
}

// Update implements repository.UserRepositoryInterface.
func (m *MockUserRepo) Update(user *entity.User) (*entity.User, error) {
	panic("unimplemented")
}

func (m *MockUserRepo) CheckSubscription(c echo.Context) (bool, error) {
	args := m.Called(c)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) Save(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepo) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestRegister(t *testing.T) {
	e := echo.New()
	mockUserRepo := new(MockUserRepo)
	cfg := &config.Config{}
	handler := NewAuthHandler(mockUserRepo, cfg)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"name":"John Doe","email":"john@example.com","password":"password"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &entity.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password",
	}
	mockUserRepo.On("Save", user).Return(user, nil)

	if assert.NoError(t, handler.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "John Doe")
	}
}

func TestLogin(t *testing.T) {
	e := echo.New()
	mockUserRepo := new(MockUserRepo)
	cfg := &config.Config{
		JWT: config.JWT{
			Secret: "secret",
		},
	}
	handler := NewAuthHandler(mockUserRepo, cfg)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"john@example.com","password":"password"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	password, _ := helpers.HashPassword("password")
	user := &entity.User{
		Email:    "john@example.com",
		Password: *password,
	}
	mockUserRepo.On("FindByEmail", "john@example.com").Return(user, nil)

	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "access_token")
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	e := echo.New()
	mockUserRepo := new(MockUserRepo)
	cfg := &config.Config{
		JWT: config.JWT{
			Secret: "secret",
		},
	}
	handler := NewAuthHandler(mockUserRepo, cfg)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"john@example.com","password":"wrongpassword"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	password, _ := helpers.HashPassword("password")
	user := &entity.User{
		Email:    "john@example.com",
		Password: *password,
	}
	mockUserRepo.On("FindByEmail", "john@example.com").Return(user, nil)

	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid credentials")
	}
}

func TestRegisterInternalServerError(t *testing.T) {
	e := echo.New()
	mockUserRepo := new(MockUserRepo)
	cfg := &config.Config{}
	handler := NewAuthHandler(mockUserRepo, cfg)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"name":"John Doe","email":"john@example.com","password":"password"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &entity.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password",
	}
	mockUserRepo.On("Save", user).Return(&entity.User{}, errors.New("internal server error"))

	if assert.NoError(t, handler.Register(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}
