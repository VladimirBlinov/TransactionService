package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/service"
	"github.com/streadway/amqp"
)

func (h *Handler) handleTransactionProcess() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &service.InputTransaction{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		body, err := json.Marshal(req)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
		}

		userQueue, err := h.rmq.Channel.QueueDeclare(
			fmt.Sprintf("user.%d", req.UserID), // name
			true,                               // durable
			false,                              // delete when unused
			false,                              // exclusive
			false,                              // no-wait
			nil,                                // arguments
		)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = h.rmq.Channel.Publish(
			"",             // exchange
			userQueue.Name, // routing key
			false,          // mandatory
			false,          // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		activeUsersQueue, err := h.rmq.Channel.QueueDeclare(
			"active_users", // name
			true,           // durable
			false,          // delete when unused
			false,          // exclusive
			false,          // no-wait
			nil,            // arguments
		)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = h.rmq.Channel.Publish(
			"",                    // exchange
			activeUsersQueue.Name, // routing key
			false,                 // mandatory
			false,                 // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		h.respond(w, r, http.StatusCreated, nil)
	}
}
