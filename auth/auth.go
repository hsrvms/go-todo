package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hsrvms/todoapp/store"
	"github.com/hsrvms/todoapp/types"
	"github.com/hsrvms/todoapp/utils"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the request (Auth header)
		tokenString := GetTokenFromRequest(r)
		// Validate the token
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Println("failed to authenticate token")
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("failed to authenticate token")
			permissionDenied(w)
			return
		}
		// Get the userId from the token
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(string)

		_, err = store.GetUserByID(userID)
		if err != nil {
			log.Println("failed to get user")
			permissionDenied(w)
			return
		}

		// Call the handler fun and continue to the endpoint
		handlerFunc(w, r)
	}
}

// CreateJWT creates a JWT token with the given secret and userID.
func CreateJWT(secret []byte, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 3).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if authHeader != "" {
		return authHeader
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func validateJWT(ts string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(ts, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteJSON(w, http.StatusUnauthorized, types.ErrorResponse{Error: "Permission denied"})
}
