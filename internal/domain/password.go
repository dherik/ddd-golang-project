package domain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	PlainPassword string //FIXME remove it, it is used only when creating a new user
	HashPassword  string
}

func (p *Password) Validate() error {

	password := p.PlainPassword

	if password == "" || len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if !hasSpecialCharacter(password) {
		return errors.New("password must contain at least one special character")
	}

	if !containsNumber(password) {
		return errors.New("password must contain at least one number")
	}

	if !containsUppercaseAndLowercase(password) {
		return errors.New("password must contain at least one uppercase letter and one lowercase letter")
	}

	return nil
}

func (p *Password) CheckPasswordHash(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.HashPassword), []byte(plainPassword))
	return err == nil
}

func (p *Password) hashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.PlainPassword), bcrypt.DefaultCost)
	p.HashPassword = string(hashedPassword)
	return err
}
