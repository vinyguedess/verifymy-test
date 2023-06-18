package handlers

import (
	"encoding/json"
	"net/http"

	"verifymy-golang-test/entities"
	"verifymy-golang-test/models"
	"verifymy-golang-test/services"
)

type signUpHandler struct {
	authService services.AuthService
}

func NewSignUpHandler(authService services.AuthService) Handler {
	return &signUpHandler{
		authService: authService,
	}
}

func (h *signUpHandler) Method() []string {
	return []string{
		http.MethodPost,
	}
}

func (h *signUpHandler) Route() string {
	return "/auth/sign_up"
}

func (h *signUpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload models.User
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = entities.NewError("Invalid JSON", []string{err.Error()})

		jsonPayload, _ := json.Marshal(err)
		w.Write(jsonPayload)
		return
	}

	credentials, err := h.authService.SignUp(r.Context(), payload)
	if err != nil {
		if _, ok := err.(*entities.EmailAlreadyInUseError); ok {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			err = entities.NewUnexpectedError(err)
		}

		jsonPayload, _ := json.Marshal(err)
		w.Write(jsonPayload)
		return
	}

	jsonPayload, _ := json.Marshal(credentials)
	w.Write(jsonPayload)
}
