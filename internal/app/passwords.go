package app

import "golang.org/x/crypto/bcrypt"

var passwordCache = map[string]bool{}

func GetPasswordCache() map[string]bool {
	return passwordCache
}

func ClearPasswordCache() {
	passwordCache = map[string]bool{}
}

func CreatePasswordHash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	return string(hashedBytes), err
}

func VerifyPasswordHash(password string, hash string) bool {
	if cached, exists := passwordCache[hash]; exists {
		return cached
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		passwordCache[hash] = false
		return false
	}

	passwordCache[hash] = true
	return true
}
