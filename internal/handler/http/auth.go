package http

import (
	"account-service/internal/domain/grant"
	"account-service/internal/domain/users"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"account-service/internal/service/auth"
	"account-service/pkg/server/response"
)

type Auth struct {
	authService *auth.Service
}

func NewAuthHandler(authService *auth.Service) *Auth {
	return &Auth{
		authService: authService,
	}
}

func (h *Auth) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/sign-up", h.signUp)
	r.Post("/sign-in", h.signIn)

	r.Route("/deps", func(r chi.Router) {
		r.Get("/", h.getDependencies)
		r.Post("/", h.setDependencies)
	})

	r.Post("/t", h.setTracking)

	return r
}

func (h *Auth) setTracking(w http.ResponseWriter, r *http.Request) {
	req := users.SobrietyTracking{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	userId, err := extractUserIDFromJWT(r)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}
	req.UserID = userId

	err = h.authService.CreateSobrietyTracking(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, response.Object{Success: true})
}

func (h *Auth) signUp(w http.ResponseWriter, r *http.Request) {
	req := grant.Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	accessToken, err := h.authService.SignUp(r.Context(), req)
	switch {
	case errors.Is(err, grant.ErrUserExist):
		response.BadRequest(w, r, err, nil)
	case errors.Is(err, nil):
		response.OK(w, r, accessToken)
	default:
		response.InternalServerError(w, r, err)
	}
	return
}

func (h *Auth) signIn(w http.ResponseWriter, r *http.Request) {
	req := grant.Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	accessToken, err := h.authService.SignIn(r.Context(), req)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		response.BadRequest(w, r, err, nil)
	case errors.Is(err, nil):
		response.OK(w, r, accessToken)
	default:
		response.InternalServerError(w, r, err)
	}
	return
}

func (h *Auth) getDependencies(w http.ResponseWriter, r *http.Request) {
	res, err := h.authService.GetDependencies(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, res)
}

func (h *Auth) setDependencies(w http.ResponseWriter, r *http.Request) {
	req := []users.UserDependency{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	userId, err := extractUserIDFromJWT(r)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	err = h.authService.CreateUserDependencies(r.Context(), userId, req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, response.Object{Success: true})
}
