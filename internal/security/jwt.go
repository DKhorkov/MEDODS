package security

import (
	"time"

	customerrors "github.com/DKhorkov/medods/internal/errors"
	"github.com/golang-jwt/jwt"
)

type JWTData struct {
	IP        string
	GUID      string
	Value     string
	SecretKey string
	Algorithm string
	TTL       time.Duration
}

func GenerateJWT(data JWTData) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(data.Algorithm))
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", customerrors.JWTClaimsError{}
	}

	claims["GUID"] = data.GUID
	claims["IP"] = data.IP
	claims["Value"] = data.Value
	claims["exp"] = time.Now().Add(data.TTL).Unix()
	return token.SignedString([]byte(data.SecretKey))
}

func ParseJWT(tokenString, secretKey string) (*JWTData, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, customerrors.InvalidJWTError{}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, customerrors.JWTClaimsError{}
	}

	data := &JWTData{
		GUID:  claims["GUID"].(string),
		IP:    claims["IP"].(string),
		Value: claims["Value"].(string),
	}

	return data, nil
}
