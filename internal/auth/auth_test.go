package auth

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestMakeJWT checks to see that we can generate a JWT and get the correct uuid out of it
func TestMakeJWT(t *testing.T) {
	uid := uuid.New()
	
	tokenStr, err := MakeJWT(uid, "test", time.Duration(10 * time.Minute))
	if err != nil {
		t.Fatalf("Failed making JWT: %v", err)
	}


	u, e := ValidateJWT(tokenStr, "test")
	if e != nil {
		t.Log(tokenStr)
		t.Log(e)
		t.Fatal("Error valdiating JWT")
	}

	if u != uid {
		t.Fatal("uuid and validatedJWT uuid does not match")
	}
}

// TestExpiredJWT makes sure that an expired jwt token errors out
func TestExpiredJWT(t *testing.T) {
	uid := uuid.New()

	tokenStr, err := MakeJWT(uid, "test", time.Duration(0 * time.Second))
	if err != nil {
		t.Fatal(fmt.Sprintf("Error making JWT: %v", err))
	}

	_, e := ValidateJWT(tokenStr, "test")
	if !strings.Contains(e.Error(), "expired") {
		t.Log(e.Error())
		t.Fatal("Token not rejected as expired, but expected to be epxired")
	}
}

// TestWrongSignature should throw an error when trying to validate with the wrong secret
func TestWrongSignature(t *testing.T) {
	uid := uuid.New()

	tokenStr, err := MakeJWT(uid, "test", time.Duration(10 * time.Minute))
	if err != nil {
		t.Fatalf("Failed creating jwt: %v", err)
	}

	_, e := ValidateJWT(tokenStr, "nope")
	if !strings.Contains(e.Error(), "token signature is invalid") {
		t.Fatal(e)
	}
}

// TestAuthHeaderExists Makes sure that GetBearerToken sees the authorization header
func TestAuthHeaderExists(t *testing.T) {
	headers := http.Header{}

	tokenStr, err := MakeJWT(uuid.New(), "test", time.Duration(10 * time.Minute))
	if err != nil {
		t.Fatalf("Making jwt failed: %v", err)
	}

	headers["Authorization"] = []string{fmt.Sprintf("Bearer %s", tokenStr)}
	tok, e := GetBearerToken(headers)
	if e != nil {
		t.Fatal(e)
	}

	if tok != tokenStr {
		t.Log(tokenStr)
		t.Log(tok)
		t.Fatal("Bearer Token did not match token string set in header")
	}
}

// TestAuthHeaderDoesNotExist
func TestAuthHeaderDoesNotExist(t *testing.T) {
	headers := http.Header{}

	_, e := GetBearerToken(headers)
	if !strings.Contains(e.Error(), "Authorization header not present") {
		t.Fatal(e)
	}
}
