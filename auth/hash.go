package auth

import "golang.org/x/crypto/bcrypt"

// Given a plaintext password p this function uses bcrypt to hash(salt is implemented by bcrypt)
// returns the hashed bytes as a string
func HashPassword(p string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(h), nil
}

// Given a hash string h and a plaintext password p,
// This function compares them and returns nil if it matched
// otherwise returns an error
func CompareHashedPasswords(h string, p string) error {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
}
