package crypt

import "testing"

func TestHashPass(t *testing.T) {
	pass := "password"
	hash, err := HashPass(pass)
	if err != nil {
		t.Error(err)
	}
	err = ComparePass(hash, pass)
	if err != nil {
		t.Error(err)
	}
}

func TestComparePass(t *testing.T) {
	pass := "password"
	hash, err := HashPass(pass)
	if err != nil {
		t.Error(err)
	}
	err = ComparePass(hash, "wrongpassword")
	if err == nil {
		t.Error("expected error, got nil")
	}
}
