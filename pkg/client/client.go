package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	defaultPageSize = 50
	tokenHeader     = "tokenId"
)

// Client is a Devolutions Server REST API client that uses Application Identity
// authentication (appKey + appSecret). Tokens expire every 5 minutes and are
// automatically refreshed.
type Client struct {
	httpClient *http.Client
	baseURL    string
	appKey     string
	appSecret  string

	mu    sync.Mutex
	token string
}

// NewClient creates a new Devolutions Server client and authenticates.
func NewClient(ctx context.Context, baseURL, appKey, appSecret string) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
		appKey:     appKey,
		appSecret:  appSecret,
	}

	if err := c.login(ctx); err != nil {
		return nil, fmt.Errorf("baton-devolutions: failed to authenticate: %w", err)
	}

	return c, nil
}

type loginResponse struct {
	Data struct {
		TokenID string `json:"tokenId"`
	} `json:"data"`
	Result  int    `json:"result"`
	Message string `json:"message"`
}

func (c *Client) login(ctx context.Context) error {
	form := url.Values{}
	form.Set("AppKey", c.appKey)
	form.Set("AppSecret", c.appSecret)

	reqURL := fmt.Sprintf("%s/api/v1/login", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("baton-devolutions: failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return fmt.Errorf("baton-devolutions: login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("baton-devolutions: login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("baton-devolutions: failed to decode login response: %w", err)
	}

	if loginResp.Data.TokenID == "" {
		return fmt.Errorf("baton-devolutions: login returned empty token")
	}

	c.mu.Lock()
	c.token = loginResp.Data.TokenID
	c.mu.Unlock()

	return nil
}

func (c *Client) ensureAuthenticated(ctx context.Context) error {
	reqURL := fmt.Sprintf("%s/api/is-logged", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return err
	}

	c.mu.Lock()
	req.Header.Set(tokenHeader, c.token)
	c.mu.Unlock()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return c.login(ctx)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.login(ctx)
	}

	// The is-logged endpoint returns a boolean.
	if string(bytes.TrimSpace(body)) != "true" {
		return c.login(ctx)
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return fmt.Errorf("baton-devolutions: authentication failed: %w", err)
	}

	reqURL := fmt.Sprintf("%s%s", c.baseURL, path)

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("baton-devolutions: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("baton-devolutions: failed to create request: %w", err)
	}

	c.mu.Lock()
	req.Header.Set(tokenHeader, c.token)
	c.mu.Unlock()

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return fmt.Errorf("baton-devolutions: request to %s failed: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("baton-devolutions: request to %s returned status %d: %s", path, resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("baton-devolutions: failed to decode response from %s: %w", path, err)
		}
	}

	return nil
}

// PaginatedResponse wraps the standard DVLS paginated API response.
type PaginatedResponse[T any] struct {
	Data        []T `json:"data"`
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalCount  int `json:"totalCount"`
	TotalPage   int `json:"totalPage"`
}

// User represents a DVLS user.
type User struct {
	ID                 string `json:"id"`
	Username           string `json:"username"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Email              string `json:"email"`
	UserType           string `json:"userType"`
	AuthenticationType string `json:"authenticationType"`
	IsEnabled          bool   `json:"isEnabled"`
	IsAdministrator    bool   `json:"isAdministrator"`
	Tags               string `json:"tags"`
	Audit              *Audit `json:"audit"`
}

// UserGroup represents a DVLS user group.
type UserGroup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Audit       *Audit `json:"audit"`
}

// GroupMember represents a user's membership in a group.
type GroupMember struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
}

// Vault represents a DVLS vault.
type Vault struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// VaultAccess represents a user or group's access to a vault.
type VaultAccess struct {
	UserID        string `json:"userId"`
	GroupID       string `json:"groupId"`
	PermissionSet string `json:"permissionSet"`
	Username      string `json:"username"`
	GroupName     string `json:"groupName"`
}

// Role represents a DVLS permission set / role.
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Audit contains audit trail timestamps.
type Audit struct {
	CreatedDate  string `json:"createdDate"`
	ModifiedDate string `json:"modifiedDate"`
}

// ListUsers returns a page of users from DVLS.
func (c *Client) ListUsers(ctx context.Context, pageNumber, pageSize int) (*PaginatedResponse[User], error) {
	path := fmt.Sprintf("/api/v3/users?pageNumber=%d&pageSize=%d", pageNumber, pageSize)
	var resp PaginatedResponse[User]
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("baton-devolutions: failed to list users: %w", err)
	}
	return &resp, nil
}

// ListGroups returns a page of user groups from DVLS.
func (c *Client) ListGroups(ctx context.Context, pageNumber, pageSize int) (*PaginatedResponse[UserGroup], error) {
	path := fmt.Sprintf("/api/v3/user-groups?pageNumber=%d&pageSize=%d", pageNumber, pageSize)
	var resp PaginatedResponse[UserGroup]
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("baton-devolutions: failed to list groups: %w", err)
	}
	return &resp, nil
}

// GetGroupMembers returns the members of a user group.
func (c *Client) GetGroupMembers(ctx context.Context, groupID string) ([]GroupMember, error) {
	path := fmt.Sprintf("/api/v3/user-groups/%s/members", groupID)
	var resp struct {
		Data []GroupMember `json:"data"`
	}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("baton-devolutions: failed to get group members for %s: %w", groupID, err)
	}
	return resp.Data, nil
}

// ListVaults returns a page of vaults from DVLS.
func (c *Client) ListVaults(ctx context.Context, pageNumber, pageSize int) (*PaginatedResponse[Vault], error) {
	path := fmt.Sprintf("/api/v3/vaults?pageNumber=%d&pageSize=%d", pageNumber, pageSize)
	var resp PaginatedResponse[Vault]
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("baton-devolutions: failed to list vaults: %w", err)
	}
	return &resp, nil
}

// GetVaultAccess returns the user and group access entries for a vault.
func (c *Client) GetVaultAccess(ctx context.Context, vaultID string) ([]VaultAccess, error) {
	path := fmt.Sprintf("/api/v3/vaults/%s/access", vaultID)
	var resp struct {
		Data []VaultAccess `json:"data"`
	}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("baton-devolutions: failed to get vault access for %s: %w", vaultID, err)
	}
	return resp.Data, nil
}

// ListRoles returns the available roles/permission sets from DVLS.
func (c *Client) ListRoles(ctx context.Context) ([]Role, error) {
	path := "/api/v3/roles"
	var resp struct {
		Data []Role `json:"data"`
	}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("baton-devolutions: failed to list roles: %w", err)
	}
	return resp.Data, nil
}

// Validate checks that the client can authenticate and make API calls.
func (c *Client) Validate(ctx context.Context) error {
	return c.ensureAuthenticated(ctx)
}
