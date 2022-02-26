package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	EmailExists(string) bool
	FindUserById(int) *User
	UserSignup(string, string)
	FindUserByEmail(string) *User
}

type UserService struct {
	store UserStore
}

func NewService(s UserStore) *UserService {
	return &UserService{
		store: s,
	}
}

func (s *UserService) Exists(email string) bool {
	return s.store.EmailExists(email)
}

func (s *UserService) FindById(id int) *User {
	return s.store.FindUserById(id)
}

func (s *UserService) SignUp(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	s.store.UserSignup(email, string(hashedPassword))

	return nil
}

func (s *UserService) SignIn(email, password string) (*User, error) {
	u := s.store.FindUserByEmail(email)
	if u.Id == 0 {
		return nil, errors.New("User does not exist")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.New("Incorrect password")
	}

	return u, nil
}