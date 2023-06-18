package handlers

import (
	"encoding/json"
	"net/http"

	"verifymy-golang-test/entities"
	"verifymy-golang-test/services"
)

type signInHandler struct {
	authService services.AuthService
}

func NewSignInHandler(authService services.AuthService) Handler {
	return &signInHandler{authService: authService}
}

func (h *signInHandler) Method() []string {
	return []string{http.MethodPost}
}

func (h *signInHandler) Route() string {
	return "/auth/sign_in"
}

func (h *signInHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = entities.NewError("Invalid JSON", []string{err.Error()})

		jsonPayload, _ := json.Marshal(err)
		w.Write(jsonPayload)

		return
	}

	credentials, err := h.authService.SignIn(
		r.Context(), payload["email"], payload["password"],
	)
	if err != nil {
		if _, ok := err.(*entities.InvalidEmailAndOrPasswordError); ok {
			w.WriteHeader(http.StatusBadRequest)
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
