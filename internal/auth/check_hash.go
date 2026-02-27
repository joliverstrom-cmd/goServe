package auth

import "github.com/alexedwards/argon2id"

func CheckPasswordHash(password string, hash string) (bool, error) {
	match, _, err := argon2id.CheckHash(password, hash)
	if err != nil {
		return false, err
	}

	return match, nil

}
