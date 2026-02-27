package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestBasicJWT(t *testing.T) {

	id := uuid.New()

	madeString, err := MakeJWT(id, "cat", time.Duration(5*time.Second))
	if err != nil {
		t.Errorf("Couldn't makeJWT: %v", err)
	}

	gotID, err := ValidateJWT(madeString, "cat")
	if err != nil {
		t.Errorf("Couldn't verify JWT: %v", err)
	}

	if id != gotID {
		t.Errorf("IDs are not matching!, ID1: %v– ID2: %v", id, gotID)
	}
}

func TestExpiredJWT(t *testing.T) {

	id := uuid.New()

	madeString, err := MakeJWT(id, "cat", time.Duration(5*time.Millisecond))
	if err != nil {
		t.Errorf("Couldn't makeJWT: %v", err)
	}

	time.Sleep(1 * time.Second)

	gotID, err := ValidateJWT(madeString, "cat")
	if err != nil {
		t.Logf("Couldn't verify JWT: %v", err)
	}

	if id == gotID {
		t.Errorf("Shouldn't pass as the token is expired")
	}
}

func TestWrongSecretJWT(t *testing.T) {

	id := uuid.New()

	madeString, err := MakeJWT(id, "cat", time.Duration(5*time.Second))
	if err != nil {
		t.Errorf("Couldn't makeJWT: %v", err)
	}

	time.Sleep(1 * time.Second)

	gotID, err := ValidateJWT(madeString, "dog")
	if err != nil {
		t.Logf("Couldn't verify JWT: %v", err)
	}

	if id == gotID {
		t.Errorf("Shouldn't pass as the token is expired")
	}
}

func TestBearerToken(t *testing.T) {
	wantedString := "WantedString"
	myHeader := http.Header{}
	myHeader.Add("Authorization", "Bearer "+wantedString)

	gotstring, err := GetBearerToken(myHeader)
	if err != nil {
		t.Errorf("Failed to get bearer: %v", err)
	}

	if gotstring != "WantedString" {
		t.Errorf("Received string: %v, wanted string: %v", gotstring, wantedString)
	}
}

func TestEmptyBearerToken(t *testing.T) {
	wantedString := "WantedString"
	myHeader := http.Header{}

	gotstring, err := GetBearerToken(myHeader)
	if err != nil {
		t.Logf("Failed to get bearer: %v", err)
		return
	}

	if gotstring == "WantedString" {
		t.Errorf("Received string: %v, wanted string: %v", gotstring, wantedString)
	}
}
