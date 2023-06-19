package handlers

import (
	"encoding/json"
	"net/http"
	"verifymy-golang-test/common"
	"verifymy-golang-test/models"
)

type showProfileHandler struct {
}

func NewShowProfileHandler() Handler {
	return &showProfileHandler{}
}

func (h *showProfileHandler) Method() []string {
	return []string{"GET"}
}

func (h *showProfileHandler) Route() string {
	return "/profile"
}

func (h *showProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(common.AuthUser).(*models.User)

	jsonPayload, _ := json.Marshal(user)
	w.Write(jsonPayload)
}
