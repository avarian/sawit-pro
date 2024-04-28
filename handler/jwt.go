package handler

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func getToken(auth string) (string, error) {
	jwtToken := strings.Split(auth, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("invalid token")
	}

	return jwtToken[1], nil
}

func (s *Server) ValidateJWT(authParam string) (*jwt.Token, error) {
	auth, err := getToken(authParam)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("invalid token")
		}
		return []byte(s.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Server) GetJWTClaims(token *jwt.Token, key string) (string, error) {
	claims := token.Claims.(jwt.MapClaims)[key].(string)
	return claims, nil
}

func (s *Server) GenerateJWT(id int) (string, error) {
	exp := time.Now().Add(time.Hour * 24).Unix() 
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  fmt.Sprint(id),
		"exp": exp,
	})

	token, err := claims.SignedString([]byte(s.Secret))
	if err != nil {
		return token, err
	}

	return token, nil
}
