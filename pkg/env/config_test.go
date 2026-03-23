package env

import "testing"

func TestGetEnvReturnsFallbackWhenMissing(t *testing.T) {
	gotEnv := getEnv("CABIN_CHAT_TEST_KEY_THAT_SHOULD_NOT_EXIST", "fallback_value")

	if gotEnv != "fallback_value" {
		t.Fatalf("expected fallback value, got: %s", gotEnv)
	}
}

func TestLoadUsesEnvValues(t *testing.T) {
	t.Setenv("CABIN_CHAT_PORT", "2222")
	t.Setenv("CABIN_CHAT_JWT_SECRET", "TestSecret")
	t.Setenv("CABIN_CHAT_DATABASE_URL", "postgres://testuser:testpassword@localhost:5432/testdb")
	t.Setenv("CABIN_CHAT_UPLOAD_DIR", "/temp-test/uploads")

	config, error := Load()
	if error != nil {
		t.Fatalf("unexpected error loading config: %v", error)
	}

	if config.Port != "2222" {
		t.Errorf("expected Port to be '2222', got: %s", config.Port)
	}

	if config.JWTSecret != "TestSecret" {
		t.Errorf("expected JWTSecret to be 'TestSecret', got: %s", config.JWTSecret)
	}

	if config.DatabaseURL != "postgres://testuser:testpassword@localhost:5432/testdb" {
		t.Errorf("expected DatabaseURL to be 'postgres://testuser:testpassword@localhost:5432/testdb', got: %s", config.DatabaseURL)
	}

	if config.UploadDir != "/temp-test/uploads" {
		t.Errorf("expected UploadDir to be '/temp-test/uploads', got: %s", config.UploadDir)
	}
}
