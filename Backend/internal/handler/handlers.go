package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	rabbit "github.com/VladimirBlinov/TransactionService/Backend/internal/rabbitmq"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	SessionName            = "MarketPlace"
	CtxKeyUser      ctxKey = iota
	ctxKeyRequestID ctxKey = iota
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type Handler struct {
	service      *service.Service
	sessionStore sessions.Store
	Router       *mux.Router
	logger       *logrus.Logger
	rmq          *rabbit.RabbitMQ
}

func NewHandler(service *service.Service, sessionStore sessions.Store, rmq *rabbit.RabbitMQ) *Handler {
	return &Handler{
		service:      service,
		sessionStore: sessionStore,
		Router:       mux.NewRouter(),
		logger:       logrus.New(),
		rmq:          rmq,
	}
}

func (h *Handler) InitHandler() {
	api := h.Router.PathPrefix("/api/v1").Subrouter()
	api.Use(handlers.CORS(
		handlers.ExposedHeaders([]string{"Set-Cookie"}),
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "content-type", "Origin", "Accept", "X-Requested-With"}),
		handlers.AllowedMethods([]string{"POST"}),
		handlers.AllowedOrigins([]string{"http://localhost"}),
	))
	api.Use(h.setRequestID)
	api.Use(h.logRequest)
	api.HandleFunc("/register", h.handleRegister()).Methods("POST")
	api.HandleFunc("/signin", h.handleSignIn()).Methods("POST")

	private := api.PathPrefix("/private").Subrouter()
	private.Use(h.AuthenticateUser)
	private.HandleFunc("/whoami", h.handleWhoami()).Methods("GET")
	private.HandleFunc("/signout", h.handleSignOut()).Methods("GET")

	transaction := private.PathPrefix("/transaction").Subrouter()
	transaction.HandleFunc("/", h.handleTransactionProcess()).Methods("POST")

}

func (h *Handler) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (h *Handler) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
