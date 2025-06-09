package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pdridh/k-line/config"
	"github.com/pdridh/k-line/db/sqlc"
)

// Generate a jwt with id and userType as the user's claims
// Returns the token as a string.
func GenerateJWT(id string, email string, name string, userType sqlc.UserType, duration time.Duration) (string, error) {
	claims := UserClaims{
		UserID:    id,
		UserEmail: email,
		UserName:  name,
		UserType:  userType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Server().JWTSecret))
}

// Validate the tokenString jwt, returns a jwt.Token ptr which has the claims inside it.
// Also checks if the signing method is the same as the generate.
// If its invalid then its returned as an error
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedJWTSigningMethod
		}

		return []byte(config.Server().JWTSecret), nil
	})
}

// A wrapper around http.SetCookie that creates a cookie named jwt with the given token
// and sets the appropriate values for it
func SetJWTCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		MaxAge:   int(config.Server().JWTExpiration.Seconds()),
		HttpOnly: true,
		Secure:   config.Server().Env == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
}

// Given a token extracts the claims as UserClaims and returns the claims
// Returns error if extraction was not succesful.
func UserClaimsFromJWT(t *jwt.Token) (*UserClaims, error) {
	c, ok := t.Claims.(*UserClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return c, nil
}
