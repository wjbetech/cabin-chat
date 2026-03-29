package auth

import (
	"testing"
	"time"
)

func TestParseAndValidateJWTAccessTokenMalformedToken(t *testing.T) {
	secret := "test-secret"
	malformedToken := "invalid-jwt-token-slop"

	_, err := ParseAndValidateJWTAccessToken(malformedToken, secret)

	if err == nil {
		t.Fatal("expected error when parsing malformed token, but got nil")
	}
}

func TestParseAndValidateJWTAccessTokenExpiredToken(t *testing.T) {
	userID := "user-123"
	secret := "test-secret"
	expiredTTL := -1 * time.Hour
	
	tokenString, err := GenerateJWTAccessToken(userID, secret, expiredTTL)
	
	if err != nil {
		t.Fatalf("GenerateJWTAccessToken returned an error: %v", err)
	}
	
	_, err = ParseAndValidateJWTAccessToken(tokenString, secret)
	
	if err == nil {
		t.Fatal("expected ParseAndValidateJWTAccessToken to return an error for an expired token, but got nil")
	}
}

func TestGenerateAccessToken(t *testing.T) {
	userID := "user-123"
	secret := "test-secret"
	ttl := time.Hour
	
	tokenString, err := GenerateJWTAccessToken(userID, secret, ttl)
	
	if err != nil {
		t.Fatalf("GenerateAccessToken returned an error: %v", err)
	}
	
	if tokenString == "" {
		t.Fatal("expected GenerateAccessToken to return a non-empty token string, but got an empty string")
	}
	
}

func TestParseAndValidateJWTAccessTokenSuccess(t *testing.T) {
	userID := "user-123"
	secret := "test-secret"
	ttl := time.Hour
	
	tokenString, err := GenerateJWTAccessToken(userID, secret, ttl)
	
	if err != nil {
		t.Fatalf("GeneratedJWTAccessToken returned an error: %v", err)
	}
	
	claims, err := ParseAndValidateJWTAccessToken(tokenString, secret)
	
	if err != nil {
		t.Fatalf("ParseAndValidateJWTAccessToken returned an error: %v", err)
	}
	
	if claims.Subject != userID {
		t.Fatalf("expected claims.Subject to be %s, got %s", userID, claims.Subject)
	}
	
	if claims.Issuer != "cabin-chat" {
		t.Fatalf("expected claims.Issuer to be 'cabin-chat', got %s", claims.Issuer)
	}
	
}

func TestParseAndValidateJWTAccessToken_InvalidSignature(t *testing.T) {
	userID := "user-123"
	secret := "test-secret"
	ttl := time.Hour
	
	tokenString, err := GenerateJWTAccessToken(userID, secret, ttl)
	
	if err != nil {
		t.Fatalf("GenerateJWTAccessToken returned an error: %v", err)
	}
	
	tamperedToken := tokenString[:len(tokenString)-1] + "x"
	
	_, err = ParseAndValidateJWTAccessToken(tamperedToken, secret)
	
	if err == nil {
		t.Fatal("expected ParseAndValidateJWTAccessToken to return an error for a token with an invalid signature, but got nil")
	}
}

