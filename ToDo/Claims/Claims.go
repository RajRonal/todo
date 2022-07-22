package Claims

import "github.com/golang-jwt/jwt"

type Claims struct {
	SessionID string `json:"session_id"`
	ID        string `json:"id"`
	jwt.StandardClaims
}
