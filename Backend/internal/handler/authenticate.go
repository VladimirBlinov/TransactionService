package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/service"
)

func (h *Handler) handleSignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &service.InputUser{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := h.service.AuthService.SignIn(req)
		if err != nil {
			h.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := h.sessionStore.Get(r, SessionName)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := h.sessionStore.Save(r, w, session); err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		h.respond(w, r, http.StatusOK, nil)
	}
}

func (h *Handler) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &service.InputUser{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := h.service.AuthService.Register(req)
		if err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		_, err = h.service.BalanceService.CreateBalance(u)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		h.respond(w, r, http.StatusCreated, u)
	}
}

func (h *Handler) handleSignOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessionStore.Get(r, SessionName)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		delete(session.Values, "user_id")
		_ = session.Save(r, w)

		h.respond(w, r.WithContext(context.Background()), http.StatusOK, nil)
	}
}

func (h *Handler) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.respond(w, r, http.StatusOK, r.Context().Value(CtxKeyUser).(*model.User))
	}
}
