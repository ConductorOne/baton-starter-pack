package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func TestNewClient_LoginSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResponse{
			Data: struct {
				TokenID string `json:"tokenId"`
			}{TokenID: "test-token-123"},
			Result: 0,
		})
	})

	server := newTestServer(t, mux)
	defer server.Close()

	client, err := NewClient(context.Background(), server.URL, "test-key", "test-secret")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client.token != "test-token-123" {
		t.Errorf("expected token 'test-token-123', got '%s'", client.token)
	}
}

func TestNewClient_LoginFailure(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResponse{
			Data: struct {
				TokenID string `json:"tokenId"`
			}{TokenID: ""},
			Result:  1,
			Message: "Invalid credentials",
		})
	})

	server := newTestServer(t, mux)
	defer server.Close()

	_, err := NewClient(context.Background(), server.URL, "bad-key", "bad-secret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListUsers_Pagination(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResponse{
			Data: struct {
				TokenID string `json:"tokenId"`
			}{TokenID: "test-token"},
		})
	})
	mux.HandleFunc("/api/is-logged", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "true")
	})
	mux.HandleFunc("/api/v3/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse[User]{
			Data: []User{
				{
					ID:        "user-1",
					Username:  "jdoe",
					FirstName: "John",
					LastName:  "Doe",
					Email:     "jdoe@example.com",
					UserType:  "User",
					IsEnabled: true,
				},
				{
					ID:        "user-2",
					Username:  "jsmith",
					FirstName: "Jane",
					LastName:  "Smith",
					Email:     "jsmith@example.com",
					UserType:  "Admin",
					IsEnabled: true,
				},
			},
			CurrentPage: 1,
			PageSize:    50,
			TotalCount:  2,
			TotalPage:   1,
		})
	})

	server := newTestServer(t, mux)
	defer server.Close()

	client, err := NewClient(context.Background(), server.URL, "key", "secret")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.ListUsers(context.Background(), 1, 50)
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}

	if len(resp.Data) != 2 {
		t.Errorf("expected 2 users, got %d", len(resp.Data))
	}

	if resp.Data[0].Username != "jdoe" {
		t.Errorf("expected username 'jdoe', got '%s'", resp.Data[0].Username)
	}

	if resp.TotalPage != 1 {
		t.Errorf("expected 1 total page, got %d", resp.TotalPage)
	}
}

func TestListGroups(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResponse{
			Data: struct {
				TokenID string `json:"tokenId"`
			}{TokenID: "test-token"},
		})
	})
	mux.HandleFunc("/api/is-logged", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "true")
	})
	mux.HandleFunc("/api/v3/user-groups", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse[UserGroup]{
			Data: []UserGroup{
				{
					ID:          "group-1",
					Name:        "Admins",
					Description: "Administrator group",
				},
			},
			CurrentPage: 1,
			PageSize:    50,
			TotalCount:  1,
			TotalPage:   1,
		})
	})

	server := newTestServer(t, mux)
	defer server.Close()

	client, err := NewClient(context.Background(), server.URL, "key", "secret")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.ListGroups(context.Background(), 1, 50)
	if err != nil {
		t.Fatalf("failed to list groups: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Errorf("expected 1 group, got %d", len(resp.Data))
	}

	if resp.Data[0].Name != "Admins" {
		t.Errorf("expected group name 'Admins', got '%s'", resp.Data[0].Name)
	}
}

func TestGetGroupMembers(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResponse{
			Data: struct {
				TokenID string `json:"tokenId"`
			}{TokenID: "test-token"},
		})
	})
	mux.HandleFunc("/api/is-logged", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "true")
	})
	mux.HandleFunc("/api/v3/user-groups/group-1/members", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Data []GroupMember `json:"data"`
		}{
			Data: []GroupMember{
				{UserID: "user-1", Username: "jdoe"},
				{UserID: "user-2", Username: "jsmith"},
			},
		})
	})

	server := newTestServer(t, mux)
	defer server.Close()

	client, err := NewClient(context.Background(), server.URL, "key", "secret")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	members, err := client.GetGroupMembers(context.Background(), "group-1")
	if err != nil {
		t.Fatalf("failed to get group members: %v", err)
	}

	if len(members) != 2 {
		t.Errorf("expected 2 members, got %d", len(members))
	}
}

func TestTokenRefresh(t *testing.T) {
	loginCount := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		loginCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResponse{
			Data: struct {
				TokenID string `json:"tokenId"`
			}{TokenID: fmt.Sprintf("token-%d", loginCount)},
		})
	})
	mux.HandleFunc("/api/is-logged", func(w http.ResponseWriter, r *http.Request) {
		// Simulate expired token on first check.
		if r.Header.Get(tokenHeader) == "token-1" {
			fmt.Fprint(w, "false")
			return
		}
		fmt.Fprint(w, "true")
	})
	mux.HandleFunc("/api/v3/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse[User]{
			Data:        []User{},
			CurrentPage: 1,
			PageSize:    50,
			TotalCount:  0,
			TotalPage:   1,
		})
	})

	server := newTestServer(t, mux)
	defer server.Close()

	client, err := NewClient(context.Background(), server.URL, "key", "secret")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// First call should trigger re-auth because is-logged returns false for token-1.
	_, err = client.ListUsers(context.Background(), 1, 50)
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}

	// login should have been called twice (initial + refresh).
	if loginCount != 2 {
		t.Errorf("expected 2 login calls, got %d", loginCount)
	}
}
