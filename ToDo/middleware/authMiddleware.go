package middleware

import (
	"ToDo/Claims"
	"ToDo/database/helper"
	"ToDo/handlers"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tkn := r.Header.Get("Authorization")

		claims := &Claims.Claims{}
		token, err := jwt.ParseWithClaims(tkn, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return handlers.JwtKey, nil
		})
		if err != nil {
			fmt.Println(err)
		}

		if token.Valid {
			ctx := context.WithValue(r.Context(), "claims", claims)
			context := ctx.Value("claims").(*Claims.Claims)
			isSession, errs := helper.SessionExist(context.SessionID)
			if errs != nil {
				log.Printf("Session Exist : Session does not exist")
				return
			}
			if !isSession {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			_, err = w.Write([]byte("Unauthorized"))
			if err != nil {
				return
			}
		}

	})
}
