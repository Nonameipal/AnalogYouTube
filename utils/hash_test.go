package utils

import "testing"

func TestGenerateHashAndCompareHash(t *testing.T) {
	hash, err := GenerateHash("secret")
	if err != nil {
		t.Fatalf("GenerateHash returned error: %v", err)
	}
	if hash == "" || hash == "secret" {
		t.Fatalf("unexpected hash value: %q", hash)
	}
	if err := CompareHash(hash, "secret"); err != nil {
		t.Fatalf("CompareHash should accept original password: %v", err)
	}
	if err := CompareHash(hash, "bad"); err == nil {
		t.Fatal("CompareHash should reject wrong password")
	}
}
