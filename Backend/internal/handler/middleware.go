package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (h *Handler) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Add("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (h *Handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := h.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
			"method":      r.Method,
			"user_id":     r.Context().Value(CtxKeyUser),
			"url":         r.URL.Path,
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof("completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start))
	})
}

func (h *Handler) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessionStore.Get(r, SessionName)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			h.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := h.service.AuthService.Authenticate(id.(int))
		if err != nil {
			h.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, u)))
	})
}
