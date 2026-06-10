package auth

import (
	"testing"
	"time"
)

func TestJWTGenerateAndVerify(t *testing.T) {
	m := NewJWTManager("test-secret", time.Hour)

	token, err := m.Generate(42, "admin")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	identity, err := m.Verify(token)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if identity.UserID != 42 || identity.Role != "admin" {
		t.Fatalf("expected userID 42 with admin role, got %+v", identity)
	}
}

func TestJWTRejectsExpiredToken(t *testing.T) {
	m := NewJWTManager("test-secret", -time.Minute)

	token, err := m.Generate(1, "user")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if _, err := m.Verify(token); err == nil {
		t.Fatal("expected error verifying expired token, got nil")
	}
}

func TestJWTRejectsWrongSecret(t *testing.T) {
	token, err := NewJWTManager("secret-a", time.Hour).Generate(1, "user")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if _, err := NewJWTManager("secret-b", time.Hour).Verify(token); err == nil {
		t.Fatal("expected error verifying token signed with another secret, got nil")
	}
}

func TestPasswordHashing(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "password123" {
		t.Fatal("hash must not equal the plaintext password")
	}
	if !CheckPassword(hash, "password123") {
		t.Fatal("CheckPassword should accept the correct password")
	}
	if CheckPassword(hash, "wrong-password") {
		t.Fatal("CheckPassword should reject a wrong password")
	}
}
