package http

import (
	"account-service/internal/domain/users"
	"account-service/internal/service/auth"
	"account-service/pkg/server/response"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type User struct {
	authService *auth.Service
}

func NewUserHandler(authService *auth.Service) *User {
	return &User{
		authService: authService,
	}
}

func (h *User) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/savings", h.GetUserSaving)
	r.Post("/savings", h.SaveUserSavings)

	return r
}

func (h *User) GetUserSaving(w http.ResponseWriter, r *http.Request) {
	userId, err := extractUserIDFromJWT(r)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	res, err := h.authService.GetUserSaving(r.Context(), userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.BadRequest(w, r, errors.New("Пользователь еще не сохранял"), nil)
			return
		}
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, res)
}

func (h *User) SaveUserSavings(w http.ResponseWriter, r *http.Request) {
	req := users.Savings{}
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

	err = h.authService.SaveUserSavings(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, response.Object{Success: true})
}
