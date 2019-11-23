package auth

import (
	"testing"
)

func TestAuthToken(t *testing.T) {
	tkStr, _, _ := GenToken("Shannon")
	ok, _, _ := AuthToken(tkStr)
	if !ok {
		t.Errorf("Token Invalid")
	}
}
