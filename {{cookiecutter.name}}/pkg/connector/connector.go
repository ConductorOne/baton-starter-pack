package connector

import (
	"context"
	"io"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
)

// Connector implements the {{ cookiecutter.name }} connector.
type Connector struct {
	// TODO: Add API client or other state here.
	// Example: client *Client
}

// ResourceSyncers returns a ResourceSyncer for each resource type.
//
// The three fundamental resource types are:
// 1. Users - principals that can be granted access (TRAIT_USER)
// 2. Groups - collections with "member" entitlement (TRAIT_GROUP)
// 3. Roles - permissions with "assigned" entitlement (TRAIT_ROLE)
func (c *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newUserBuilder(c),
		newGroupBuilder(c),
		newRoleBuilder(c),
	}
}

// Asset fetches an asset by reference. Most connectors return nil.
func (c *Connector) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns connector metadata shown in the UI.
func (c *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "{{ cookiecutter.name }}",
		Description: "Connector for {{ cookiecutter.name }}",
	}, nil
}

// Validate tests the connection. Called before every sync.
func (c *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	// TODO: Test API connection
	// Example:
	//   _, err := c.client.GetCurrentUser(ctx)
	//   if err != nil {
	//       return nil, fmt.Errorf("{{ cookiecutter.name }}: validation failed: %w", err)
	//   }
	return nil, nil
}

// New creates a new connector instance.
func New(ctx context.Context) (*Connector, error) {
	// TODO: Initialize API client
	return &Connector{}, nil
}
