package http

import (
	"main/config"
	"main/http/handler"
	"main/http/middleware"
	"main/repository"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Route struct {
	Method  string
	Path    string
	Handler echo.HandlerFunc
	IsAuth  bool
}

func BuildServer(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	routes := []Route{}

	// init middleware
	middlewareAuth := middleware.AuthMiddleware(cfg.JWT.Secret)

	// init repository
	userRepository := repository.NewUserRepository(db)
	profileRepository := repository.NewProfileRepository(db)
	matchRepository := repository.NewMatchRepository(db)

	// init handler
	authHandler := handler.NewAuthHandler(userRepository, cfg)
	datingHandler := handler.NewDatingHandler(profileRepository, matchRepository)
	userHandler := handler.NewUserHandler(userRepository, profileRepository)

	// init routes
	authRoutes := routeAuth(authHandler)
	datingRoutes := routeDating(datingHandler)
	profileRoutes := routeProfile(userHandler)
	routes = append(routes, (*authRoutes)...)
	routes = append(routes, (*datingRoutes)...)
	routes = append(routes, (*profileRoutes)...)
	for _, route := range routes {
		if route.IsAuth {
			e.Add(route.Method, route.Path, route.Handler, middlewareAuth)
		} else {
			e.Add(route.Method, route.Path, route.Handler)
		}
	}
}

func routeAuth(h *handler.AuthHandler) *[]Route {
	authRoutes := []Route{}
	loginRoute := Route{
		Method:  "POST",
		IsAuth:  false,
		Path:    "/login",
		Handler: h.Login,
	}

	registerRoute := Route{
		Method:  "POST",
		IsAuth:  false,
		Path:    "/register",
		Handler: h.Register,
	}

	authRoutes = append(authRoutes, loginRoute, registerRoute)
	return &authRoutes
}

func routeProfile(h *handler.UserHandler) *[]Route {
	profileRoutes := []Route{}
	meRoute := Route{
		Method:  "GET",
		IsAuth:  true,
		Path:    "/me",
		Handler: h.Me,
	}

	updateMeRoute := Route{
		Method:  "PUT",
		IsAuth:  true,
		Path:    "/me",
		Handler: h.UpdateProfile,
	}

	purchasePremiumRoute := Route{
		Method:  "POST",
		IsAuth:  true,
		Path:    "/subscribe",
		Handler: h.PurchasePremium,
	}

	profileRoutes = append(profileRoutes, meRoute, purchasePremiumRoute, updateMeRoute)
	return &profileRoutes
}

func routeDating(h *handler.DatingHandler) *[]Route {
	datingRoutes := []Route{}
	profileRoute := Route{
		Method:  "GET",
		IsAuth:  true,
		Path:    "/profile",
		Handler: h.Profile,
	}

	swipedProfileRoute := Route{
		Method:  "POST",
		IsAuth:  true,
		Path:    "/swipe",
		Handler: h.SwipedProfile,
	}

	matchRoute := Route{
		Method:  "GET",
		IsAuth:  true,
		Path:    "/match",
		Handler: h.MatchList,
	}

	datingRoutes = append(datingRoutes, profileRoute, swipedProfileRoute, matchRoute)
	return &datingRoutes
}
