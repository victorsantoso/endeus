package helper

import (
	"errors"
	"strconv"
	"time"

	"github.com/apex/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/victorsantoso/endeus/internal"
)

var (
	jwtConfig              = internal.ConfigureJWT()
	ErrInvalidJwtAlgorithm = errors.New("invalid jwt signing method")
)

func GenerateJWT(role string, userId int64) (string, error) {
	id := strconv.Itoa(int(userId))
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{ // can use custom claim if you need to inject more data
		ID:        id,
		Audience:  jwt.ClaimStrings{jwtConfig.Audience},
		Issuer:    jwtConfig.Issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
		Subject:   role,
	})
	token, err := jwtToken.SignedString([]byte(jwtConfig.Secret)) // sign token with private key
	if err != nil {
		log.Debugf("[jwt.GenerateJWT] err signing jwt: %v\n", err)
		return "", err // error signing token
	}
	// return signed token
	return token, nil
}

func VerifyJwt(jwtString string) (*jwt.Token, error) {
    jwtToken, err := jwt.Parse(jwtString, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
            return nil, ErrInvalidJwtAlgorithm
        }
        return []byte(jwtConfig.Secret), nil
    })
    if err != nil {
        return nil, err
    }
    return jwtToken, nil
}