package security

import "golang.org/x/crypto/bcrypt"

func HashRefreshToken(token string, hashCost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(token), hashCost)
	return string(bytes), err
}

func ValidateRefreshToken(token, hashedToken string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
	return err == nil
}
