package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hsrvms/todoapp/models"
	"github.com/hsrvms/todoapp/store"
)

func TestRegisterUser(t *testing.T) {
	testCases := []struct {
		name     string
		payload  *models.User
		expCode  int
		expError error
	}{
		{
			name: "empty username",
			payload: &models.User{
				Username: "",
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "empty password",
			payload: &models.User{
				Username: "testUser",
				Password: "",
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "valid user",
			payload: &models.User{
				Username: "testUserRegister",
				Password: "testPassword",
			},
			expCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms := &store.MockStore{}
			service := NewUserService(ms)

			if service == nil {
				t.Fatal("failed to create UserService")
			}

			b, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(b))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			res := httptest.NewRecorder()

			mux := http.NewServeMux()
			service.RegisterRoutes(mux, "")

			if mux == nil {
				t.Fatal("failed to create ServeMux")
			}

			mux.ServeHTTP(res, req)

			if res.Code != tc.expCode {
				t.Errorf("got %d want %d", res.Code, tc.expCode)
			}

		})
	}
}

func TestLoginUser(t *testing.T) {
	for _, tc := range []struct {
		name     string
		payload  *models.User
		expCode  int
		expError error
	}{
		{
			name: "empty username",
			payload: &models.User{
				Username: "",
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "empty password",
			payload: &models.User{
				Username: "testUser",
				Password: "",
			},
			expCode: http.StatusBadRequest,
		},
		{
			// TODO: This doesn't work so far. Need to figure out why.
			name: "valid user",
			payload: &models.User{
				Username: "testUserLogin",
				Password: "testPassword",
			},
			expCode: http.StatusCreated,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ms := store.NewMockStore()
			service := NewUserService(ms)
			if service == nil {
				t.Fatal("failed to create UserService")
			}

			payload, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(payload))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			res := httptest.NewRecorder()

			mux := http.NewServeMux()
			service.RegisterRoutes(mux, "")
			if mux == nil {
				t.Fatal("failed to create ServeMux")
			}

			mux.ServeHTTP(res, req)

			if res.Code != tc.expCode {
				t.Errorf("got %d want %d", res.Code, tc.expCode)
			}
		})
	}
}
