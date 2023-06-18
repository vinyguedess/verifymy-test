package handlers

import (
	"encoding/json"
	"net/http"
	"verifymy-golang-test/common"
	"verifymy-golang-test/models"
)

type profileHandler struct {
}

func NewProfileHandler() Handler {
	return &profileHandler{}
}

func (h *profileHandler) Method() []string {
	return []string{"GET"}
}

func (h *profileHandler) Route() string {
	return "/profile"
}

func (h *profileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(common.AuthUser).(*models.User)

	jsonPayload, _ := json.Marshal(user)
	w.Write(jsonPayload)
}
