package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-devolutions/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
)

type Devolutions struct {
	client *client.Client
}

func New(ctx context.Context, serverURL, appKey, appSecret string) (*Devolutions, error) {
	c, err := client.NewClient(ctx, serverURL, appKey, appSecret)
	if err != nil {
		return nil, fmt.Errorf("baton-devolutions: failed to create client: %w", err)
	}

	return &Devolutions{client: c}, nil
}

func (d *Devolutions) Metadata(_ context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Devolutions Server",
		Description: "Connector syncing users, groups, roles, and vaults from Devolutions Server (DVLS) to ConductorOne.",
	}, nil
}

func (d *Devolutions) Validate(ctx context.Context) (annotations.Annotations, error) {
	if err := d.client.Validate(ctx); err != nil {
		return nil, fmt.Errorf("baton-devolutions: validation failed: %w", err)
	}
	return nil, nil
}

func (d *Devolutions) ResourceSyncers(_ context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newUserBuilder(d.client),
		newGroupBuilder(d.client),
		newRoleBuilder(d.client),
		newVaultBuilder(d.client),
	}
}
