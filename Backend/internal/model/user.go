package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int     `json:"id"`
	Email             string  `json:"email"`
	Password          string  `json:"password,omitempty"`
	EncryptedPassword string  `json:"-"`
	Balance           float32 `json:"balance"`
	BalanceID         int     `json:"balance_id"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(4, 50)),
		validation.Field(&u.Balance, validation.Min(0)),
	)
}

func (u *User) EncryptPasswordBeforeCreate() error {
	if len(u.Password) > 0 {
		encryptedString, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = encryptedString
	}

	return nil
}

func (u *User) Sanitize() {
	u.Password = ""
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
