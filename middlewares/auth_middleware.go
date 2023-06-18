package middlewares

import (
	"context"
	"net/http"
	"strings"

	"verifymy-golang-test/common"
	"verifymy-golang-test/services"
	"verifymy-golang-test/utils"
)

var ALLOWED_PATHS = []interface{}{
	"/",
	"/auth/sign_in",
	"/auth/sign_up",
}

func AuthMiddleware(
	authService services.AuthService,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !utils.SliceContains(ALLOWED_PATHS, r.URL.Path) {
				authorizationHeader := strings.Split(r.Header.Get("Authorization"), " ")
				if len(authorizationHeader) != 2 {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"message": "Malformed authorization header"}`))
					return
				}

				if strings.ToLower(authorizationHeader[0]) != "bearer" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"message": "Authorization header must be a bearer token"}`))
					return
				}

				ctx := r.Context()

				accessToken := authorizationHeader[1]
				user, err := authService.GetUserFromToken(ctx, accessToken)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"message": "Invalid access token"}`))
					return
				}

				ctx = context.WithValue(ctx, common.AuthUser, user)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}
