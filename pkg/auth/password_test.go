package auth

import "testing"

func TestHashPassword(t *testing.T) {
	plainPassword := "verysecretpassword"

	hashedPassword, err := HashPassword(plainPassword)

	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if hashedPassword == "" {
		t.Fatal("expected hashed password is not expected to be empty")
	}

	if hashedPassword == plainPassword {
		t.Fatal("hashedPassword is expected to differ from plainPassword")
	}
}

func TestCheckPassword_Success(t *testing.T) {
	plainPassword := "verysecretpassword"

	hashedPassword, err := HashPassword(plainPassword)

	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)

	}

	err = CheckPassword(plainPassword, hashedPassword)

	if err != nil {
		t.Fatalf("CheckPassword returned an error: %v", err)
	}
}

func TestCheckPassword_Failure(t *testing.T) {
	plainPassword := "verysecretpassword"
	wrongPassword := "wrongpassword"

	hashedPassword, err := HashPassword(plainPassword)

	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	err = CheckPassword(wrongPassword, hashedPassword)

	if err == nil {
		t.Fatal("CheckPassword did not return an error for incorrect password")
	}
}