package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	hashedPW, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hashedPW, nil
}
