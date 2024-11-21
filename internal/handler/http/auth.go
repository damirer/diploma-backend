package http

import (
	"account-service/internal/domain/grant"
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

	return r
}

// @Summary		Авторизация
// @Description	Авторизация
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		grant.Request	true	"request body"
// @Success		200		{object}	response.Object
// @Failure		400		{object}	response.Object
// @Failure		500		{object}	response.Object
// @Router			/auth [post]
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
