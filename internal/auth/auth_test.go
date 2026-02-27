package auth

import (
	"testing"
)

func TestHash(t *testing.T) {
	pw := "vitanVitoni"
	hashedPW, _ := HashPassword(pw)

	check, _ := CheckPasswordHash(pw, hashedPW)

	if !check {
		t.Errorf(`Password %s does not match with hash: %s`, pw, hashedPW)
	}

}

func TestHashFail(t *testing.T) {
	pw := "vitanVitoni"
	checkedPW := "vitanvitoni"
	hashedPW, _ := HashPassword(pw)

	check, _ := CheckPasswordHash("vitanvitoni", hashedPW)

	if check {
		t.Errorf(`Password %s does not match with hash: %s`, checkedPW, hashedPW)
		return
	}

	t.Logf("Password: %s, checkedPW: %s do NOT match", pw, checkedPW)

}
