package tests

import (
	"testing"

	"liokoredu/pkg/generators"
)

func TestHash(t *testing.T) {
	password := "Wolf123"

	hash := generators.HashPassword(password)

	if !generators.CheckHashedPassword(hash, password) {
		t.Errorf("password failed hash check")
	}
}
