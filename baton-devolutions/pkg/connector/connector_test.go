package connector

import (
	"testing"

	"github.com/conductorone/baton-devolutions/pkg/client"
)

func TestResourceTypes(t *testing.T) {
	tests := []struct {
		name     string
		rt       string
		expected string
	}{
		{"user", resourceTypeUser.Id, "user"},
		{"group", resourceTypeGroup.Id, "group"},
		{"role", resourceTypeRole.Id, "role"},
		{"vault", resourceTypeVault.Id, "vault"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.rt != tt.expected {
				t.Errorf("expected resource type ID '%s', got '%s'", tt.expected, tt.rt)
			}
		})
	}
}

func TestDisplayName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		username  string
		expected  string
	}{
		{"full name", "John", "Doe", "jdoe", "John Doe"},
		{"first only", "John", "", "jdoe", "John"},
		{"last only", "", "Doe", "jdoe", "Doe"},
		{"username fallback", "", "", "jdoe", "jdoe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := client.User{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
				Username:  tt.username,
			}
			result := displayName(user)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
