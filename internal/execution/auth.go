package execution

import (
	"encoding/hex"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func readJwt(jwtPath string) string {
	data, err := os.ReadFile(jwtPath)
	if err != nil {
		panic(err)
	}
	jwt := strings.TrimSpace(string(data))
	jwt = strings.TrimLeft(jwt, "0x")
	return jwt
}

func loadJwt(jwtPath string) []byte {
	jwt := readJwt(jwtPath)
	data, err := hex.DecodeString(jwt)
	if err != nil {
		panic(err)
	}
	return data
}

func genToken(secret []byte) string {
	if secret == nil {
		panic("no secret")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}
	return tokenString
}
