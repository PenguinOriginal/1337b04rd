package utils

import (
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	uuid, err := GenerateUUID()
	if err != nil {
		t.Fatalf("GenerateUUID returned an error: %v", err)
	}

	if uuid == "" {
		t.Error("Expected non-empty UUID string")
	}

	// Check basic UUID format
	if len(uuid) != 36 {
		t.Errorf("Expected UUID length of 36, got %d: %s", len(uuid), uuid)
	}

	if uuid[14] != '4' {
		t.Errorf("Expected version 4 UUID, got: %s", uuid)
	}
}
