package handler

import (
	"net/http"
	"todotree/auth"
)

type httpHandler = http.Handler
type middleware = func(next httpHandler) httpHandler

func AuthMiddleware(j *auth.JWTer) middleware {
	return func(next httpHandler) httpHandler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, err := j.FillContext(r)
			if err != nil {
				RespondJSON(r.Context(), w, ErrResponse{
					Message: "not find auth info",
					Details: []string{err.Error()},
				}, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

func AdminMiddleware(next httpHandler) httpHandler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsAdmin(r.Context()) {
			RespondJSON(r.Context(), w, ErrResponse{
				Message: "not admin",
			}, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
