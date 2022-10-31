package handler

import (
	"net/http"
)

func (h *Handler) handleTransactionProcess() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//req := &service.InputProduct{}

		h.respond(w, r, http.StatusCreated, nil)
	}
}
