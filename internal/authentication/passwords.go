package authentication

import "golang.org/x/crypto/bcrypt"

// TODO: Re-add password cache

func CreatePasswordHash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(hashedBytes), err
}

func VerifyPasswordHash(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
