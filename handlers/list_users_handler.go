package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"verifymy-golang-test/entities"
	"verifymy-golang-test/services"
)

type listUsersHandler struct {
	userService services.UserService
}

func NewListUsersHandler(
	userService services.UserService,
) Handler {
	return &listUsersHandler{
		userService: userService,
	}
}

func (h *listUsersHandler) Method() []string {
	return []string{http.MethodGet}
}

func (h *listUsersHandler) Route() string {
	return "/users"
}

func (h *listUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}
	intLimit, _ := strconv.Atoi(limit)

	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	intPage, _ := strconv.Atoi(page)

	users, count, err := h.userService.FindAll(r.Context(), intLimit, intPage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = entities.NewUnexpectedError(err)

		jsonPayload, _ := json.Marshal(err)
		w.Write(jsonPayload)
		return
	}

	jsonPayload, _ := json.Marshal(users)
	w.Header().Set("X-Total-Count", strconv.FormatInt(count, 10))
	w.Write(jsonPayload)
}
