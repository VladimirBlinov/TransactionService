package apiserver

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/handler"
	rabbit "github.com/VladimirBlinov/TransactionService/Backend/internal/rabbitmq"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/service"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store/sqlstore"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type ApiServer struct {
	httpServer *http.Server
}

func (s *ApiServer) Start(config *Config) error {
	rmq, err := rabbit.NewRabbitMQ(config.RabbitURL)
	if err != nil {
		return err
	}
	defer rmq.Close()

	db, err := newDB(config.DataBaseURL)
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			logrus.Errorf("error db close: %s", err.Error())
		}
	}(db)

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	services := service.NewService(store)
	handlers := handler.NewHandler(services, sessionStore, rmq)
	handlers.InitHandler()

	s.httpServer = &http.Server{
		Addr:           config.BindAddr,
		Handler:        handlers.Router,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (s *ApiServer) ShutDown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
