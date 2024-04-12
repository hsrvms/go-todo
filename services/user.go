package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hsrvms/todoapp/models"
	"github.com/hsrvms/todoapp/store"
	"github.com/hsrvms/todoapp/utils"
)

var ErrUsernameRequired = errors.New("username is required")
var ErrPasswordRequired = errors.New("password is required")

type UserService struct {
	store store.Store
}

func NewUserService(store store.Store) *UserService {
	return &UserService{store: store}
}

// POST /auth/register:
//
// Payload:
//
//	{"username": "johnDoe", "password": "secretPassword"}
//
// POST /auth/login:
//
// Payload:
//
//	{"username": "johnDoe", "password": "secretPassword"}
func (s *UserService) RegisterRoutes(mux *http.ServeMux, prefix string) {
	endpointRegister := generateEndpoint("POST", prefix, "/auth/register")
	endpointLogin := generateEndpoint("POST", prefix, "/auth/login")

	mux.HandleFunc(endpointRegister, s.handleUserRegister)
	mux.HandleFunc(endpointLogin, s.handleUserLogin)
}

func (s *UserService) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := decodeJSON(r, &user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validateUserPayload(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPW, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPW

	createdUser, err := s.store.CreateUser(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	if createdUser == nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	token, err := createAndSetAuthCookie(createdUser.ID, w)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, token)
}

func (s *UserService) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := decodeJSON(r, &user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validateUserPayload(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingUser, err := s.store.GetUserByUsername(user.Username)
	if err != nil || existingUser == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	fmt.Printf("existingUser: %v \n newUser: %v\n", existingUser, user)
	match := checkPasswordHash(user.Password, existingUser.Password)
	if !match {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := createAndSetAuthCookie(existingUser.ID, w)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(token); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func validateUserPayload(user *models.User) error {
	if user.Username == "" {
		return ErrUsernameRequired
	}

	if user.Password == "" {
		return ErrPasswordRequired
	}

	return nil
}
