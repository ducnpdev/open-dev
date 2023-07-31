package jwttoken

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func ParserToken(tokenString string) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<YOUR VERIFICATION KEY>"), nil
	})
	// ... error handling
	fmt.Println(token, err)
	// do something with decoded claims
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)
	}
}
