package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var (
	errEmailRequired     = errors.New("email is required")
	errFirstNameRequired = errors.New("firstname is required")
	errLastNameRequired  = errors.New("lastname is required")
	errPasswordRequired  = errors.New("password is required")
)

type UserService struct {
	store Store
}

func NewUserService(s Store) *UserService {
	return &UserService{store: s}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.handleUserRegister).Methods("POST")
}

func (s *UserService) handleUserRegister(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)

	if err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	defer r.Body.Close()

	var payload *User

	err = json.Unmarshal(body, &payload)

	if err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := validateUserPayload(payload); err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	hashedPW, err := HashPassword(payload.Password)

	if err != nil {
		WriteJson(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	payload.Password = hashedPW

	u, err := s.store.CreateUser(payload)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJson(w, http.StatusCreated, token)

}

func validateUserPayload(user *User) error {
	if user.Email == "" {
		return errEmailRequired
	}

	if user.FirstName == "" {
		return errFirstNameRequired
	}

	if user.LastName == "" {
		return errLastNameRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	secret := []byte(Envs.JWTSecret)

	token, err := CreateJWT(secret, id)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}
