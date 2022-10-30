package service

import (
	"errors"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store"
)

type InputUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthService struct {
	store store.Store
}

func (s *AuthService) Register(req *InputUser) (*model.User, error) {
	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := s.store.User().Create(u); err != nil {
		return nil, err
	}

	u.Sanitize()

	return u, nil
}

func (s *AuthService) SignIn(req *InputUser) (*model.User, error) {
	u, err := s.store.User().FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if !u.ComparePassword(req.Password) {
		return nil, errors.New("invalid Password")
	}

	return u, nil
}

func (s *AuthService) Authenticate(id int) (*model.User, error) {
	u, err := s.store.User().FindById(id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func NewAuthService(store store.Store) *AuthService {
	return &AuthService{
		store: store,
	}
}
