package configs

import "testing"

func TestEnvHelpers(t *testing.T) {
	t.Setenv("STRING_VALUE", "hello")
	if got := envString("STRING_VALUE", "fallback"); got != "hello" {
		t.Fatalf("expected env value, got %q", got)
	}
	if got := envString("MISSING_STRING_VALUE", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback value, got %q", got)
	}

	t.Setenv("INT_VALUE", "42")
	if got := envInt("INT_VALUE", 1); got != 42 {
		t.Fatalf("expected parsed int, got %d", got)
	}
	t.Setenv("BAD_INT_VALUE", "oops")
	if got := envInt("BAD_INT_VALUE", 5); got != 5 {
		t.Fatalf("expected fallback int, got %d", got)
	}
}

func TestEnvRequired(t *testing.T) {
	t.Setenv("REQUIRED_VALUE", " secret ")
	got, err := envRequired("REQUIRED_VALUE")
	if err != nil {
		t.Fatalf("envRequired returned error: %v", err)
	}
	if got != "secret" {
		t.Fatalf("expected trimmed value, got %q", got)
	}

	t.Setenv("EMPTY_REQUIRED_VALUE", " ")
	if _, err := envRequired("EMPTY_REQUIRED_VALUE"); err == nil {
		t.Fatal("expected error for empty required value")
	}
}

func TestLoadRequiresDefaultDatabaseName(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("POSTGRES_USER", "postgres")
	t.Setenv("POSTGRES_PASSWORD", "password")
	t.Setenv("POSTGRES_DATABASE", "wrong")

	if err := Load(); err == nil {
		t.Fatal("expected error for wrong database name")
	}
}

func TestLoadBuildsSettingsFromEnvironment(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("POSTGRES_USER", "postgres")
	t.Setenv("POSTGRES_PASSWORD", "password")
	t.Setenv("POSTGRES_DATABASE", defaultPostgresDatabase)
	t.Setenv("SERVER_PORT", "8081")
	t.Setenv("ACCESS_TOKEN_TTL_MINUTES", "20")
	t.Setenv("REFRESH_TOKEN_TTL_DAYS", "40")
	t.Setenv("REDIS_DB", "2")

	if err := Load(); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if AppSettings.AppParams.PortRun != "8081" {
		t.Fatalf("expected port 8081, got %q", AppSettings.AppParams.PortRun)
	}
	if AppSettings.PostgresParams.Database != defaultPostgresDatabase {
		t.Fatalf("database name changed: %q", AppSettings.PostgresParams.Database)
	}
	if AppSettings.AuthParams.AccessTokenTtlMinutes != 20 || AppSettings.AuthParams.RefreshTokenTtlDays != 40 {
		t.Fatalf("unexpected token ttl settings: %+v", AppSettings.AuthParams)
	}
	if AppSettings.RedisParams.DB != 2 {
		t.Fatalf("expected redis db 2, got %d", AppSettings.RedisParams.DB)
	}
}
