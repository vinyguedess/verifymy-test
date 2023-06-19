package handlers

import (
	"encoding/json"
	"net/http"

	"verifymy-golang-test/entities"
	"verifymy-golang-test/models"
	"verifymy-golang-test/services"
)

type updateProfileHandler struct {
	userService services.UserService
}

func NewUpdateProfileHandler(
	userService services.UserService,
) Handler {
	return &updateProfileHandler{
		userService: userService,
	}
}

func (h *updateProfileHandler) Method() []string {
	return []string{http.MethodPut}
}

func (h *updateProfileHandler) Route() string {
	return "/profile"
}

func (h *updateProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload models.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusUnprocessableEntity)
		err = entities.NewError("Invalid JSON", []string{err.Error()})

		jsonPayload, _ := json.Marshal(err)
		w.Write(jsonPayload)
		return
	}

	if err := h.userService.UpdateProfile(r.Context(), payload); err != nil {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusInternalServerError)
		err = entities.NewUnexpectedError(err)

		jsonPayload, _ := json.Marshal(err)
		w.Write(jsonPayload)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
