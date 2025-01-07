package handler

import (
	"encoding/json"
	"materials/internal/config"
	"materials/internal/dto"
	"net/http"
	"time"

	"github.com/go-playground/validator"
)

type UserHandler struct {
	service UserService
	cfg     *config.Config
}

type UserService interface {
	Registration(registrationDTO dto.RegistrationUser) (*dto.RegistrationReturn, error)
	Login(loginDTO dto.LoginUser) (string, string, error)
}

func NewUserHandler(s UserService, c *config.Config) *UserHandler {
	return &UserHandler{
		service: s,
		cfg:     c,
	}
}

func (h *UserHandler) Registration(rw http.ResponseWriter, r *http.Request) {
	var user dto.RegistrationUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendJSONError(rw, "invalid input", http.StatusBadRequest)
		return
	}

	if err := validateRegistrationUser(user); err != nil {
		sendJSONError(rw, "not valid login or password", http.StatusBadRequest)
		return
	}

	regReturn, err := h.service.Registration(user)
	if err != nil {
		sendJSONError(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshTokenCokie := http.Cookie{
		Name:     "refreshtoken",
		Value:    regReturn.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(2 * time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   h.cfg.TokensSecure,
	}

	accessTokenCookie := http.Cookie{
		Name:     "accesstoken",
		Value:    regReturn.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   h.cfg.TokensSecure,
	}

	http.SetCookie(rw, &refreshTokenCokie)
	http.SetCookie(rw, &accessTokenCookie)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(rw http.ResponseWriter, r *http.Request) {
	var user dto.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendJSONError(rw, "invalid input", http.StatusBadRequest)
		return
	}

	access, refresh, err := h.service.Login(user)
	if err != nil {
		sendJSONError(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshTokenCokie := http.Cookie{
		Name:     "refreshtoken",
		Value:    refresh,
		Path:     "/",
		Expires:  time.Now().Add(2 * time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   h.cfg.TokensSecure,
	}

	accessTokenCookie := http.Cookie{
		Name:     "accesstoken",
		Value:    access,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   h.cfg.TokensSecure,
	}

	http.SetCookie(rw, &refreshTokenCokie)
	http.SetCookie(rw, &accessTokenCookie)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
}

func validateRegistrationUser(user dto.RegistrationUser) error {
	validate := validator.New()
	return validate.Struct(user)
}

func sendJSONError(rw http.ResponseWriter, message string, status int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	errResponse := map[string]string{"error": message}
	json.NewEncoder(rw).Encode(errResponse)
}
