package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"verifymy-golang-test/entities"
	"verifymy-golang-test/services"
)

type deleteUserByIdHandler struct {
	userService services.UserService
}

func NewDeleteUserByIdHandler(
	userService services.UserService,
) Handler {
	return &deleteUserByIdHandler{
		userService: userService,
	}
}

func (h *deleteUserByIdHandler) Method() []string {
	return []string{http.MethodDelete}
}

func (h *deleteUserByIdHandler) Route() string {
	return "/users/{user_id}"
}

func (h *deleteUserByIdHandler) ServeHTTP(
	w http.ResponseWriter, r *http.Request,
) {
	params := mux.Vars(r)

	userId := params["user_id"]
	user, err := h.userService.FindById(r.Context(), userId)
	if err != nil {
		if _, ok := err.(*entities.ItemNotFoundError); ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			err = entities.NewUnexpectedError(err)
		}

		jsonPayload, _ := json.Marshal(err)
		w.Write(jsonPayload)

		return
	}

	err = h.userService.DeleteById(r.Context(), user.ID.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = entities.NewUnexpectedError(err)

		jsonPayload, _ := json.Marshal(err)
		w.Write(jsonPayload)
	}

	w.WriteHeader(http.StatusNoContent)
}
