package tools

import (
	"testing"
)

func TestAuthToken(t *testing.T) {
	tkStr, e, _ := GenToken("Shannon")
	if e != nil {
		t.Error(e)
	}
	_, err, _ := AuthToken(tkStr)
	if err != nil {
		t.Error(err)
	}
}
