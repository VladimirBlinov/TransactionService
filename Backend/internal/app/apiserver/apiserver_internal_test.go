package apiserver_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/handler"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/service"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store/sqlstore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

func TestServerHandleSignOut(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "transaction", "user_transaction", "user_balance", "balance_audit", "balance")

	store := sqlstore.New(db)
	services := service.NewService(store)

	u := model.TestUser(t)
	store.User().Create(u)

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(services, sessions.NewCookieStore(secretKey))
	handlers.InitHandler()
	sc := securecookie.New(secretKey, nil)

	testCases := []struct {
		name         string
		context      *model.User
		coockieValue map[interface{}]interface{}
		expectedCode int
	}{
		{
			name:    "valid",
			context: u,
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/private/signout", nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			ctx := context.WithValue(req.Context(), handler.CtxKeyUser, tc.context)
			handlers.Router.ServeHTTP(rec, req.WithContext(ctx))
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.NotEqual(t, tc.coockieValue, rec.Result().Header["Set-Cookie"])
		})
	}
}

func TestServer_AuthenticateUser(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "transaction", "user_transaction", "user_balance", "balance_audit", "balance")

	store := sqlstore.New(db)
	srvc := service.NewService(store)

	u := model.TestUser(t)
	store.User().Create(u)

	secretKey := []byte("secret_key")
	handlers := handler.NewHandler(srvc, sessions.NewCookieStore(secretKey))
	handlers.InitHandler()
	sc := securecookie.New(secretKey, nil)
	handl := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	testCases := []struct {
		name         string
		coockieValue map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			coockieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			coockieValue: nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/", nil)
			coockieStr, _ := sc.Encode(handler.SessionName, tc.coockieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", handler.SessionName, coockieStr))
			handlers.AuthenticateUser(handl).ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleRegister(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "transaction", "user_transaction", "user_balance", "balance_audit", "balance")

	store := sqlstore.New(db)
	srvc := service.NewService(store)

	handlers := handler.NewHandler(srvc, sessions.NewCookieStore([]byte("secret_key")))
	handlers.InitHandler()
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "user@example.org",
				"password": "password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", b)
			handlers.Router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleSignIn(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "transaction", "user_transaction", "user_balance", "balance_audit", "balance")

	store := sqlstore.New(db)
	srvc := service.NewService(store)

	u := model.TestUser(t)
	store.User().Create(u)

	handlers := handler.NewHandler(srvc, sessions.NewCookieStore([]byte("secret_key")))
	handlers.InitHandler()
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    u.Email,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid",
				"password": u.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    u.Email,
				"password": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/signin", b)
			req.Header.Set("Origin", "http://localhost")
			handlers.Router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
