package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hsrvms/todoapp/auth"
	"golang.org/x/crypto/bcrypt"
)

func generateEndpoint(method, prefix, path string) string {
	return fmt.Sprintf("%v %v%v", method, prefix, path)
}

func decodeJSON(r *http.Request, v any) error {
	fmt.Println(r.Body)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.Unmarshal(body, v)
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	log.Println(err)
	return err == nil
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := auth.CreateJWT([]byte(secret), id)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}
