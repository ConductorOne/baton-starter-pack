package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-devolutions/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	sdkGrant "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const (
	vaultsPageSize = 50

	// DVLS permission sets for vault access.
	permissionContributor = "Contributor"
	permissionOperator    = "Operator"
	permissionReader      = "Reader"
)

var vaultPermissions = []string{
	permissionContributor,
	permissionOperator,
	permissionReader,
}

type vaultBuilder struct {
	resourceType *v2.ResourceType
	client       *client.Client
}

func (v *vaultBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return v.resourceType
}

func vaultResource(vault client.Vault) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"vault_name": vault.Name,
		"vault_id":   vault.ID,
	}

	if vault.Description != "" {
		profile["description"] = vault.Description
	}

	ret, err := rs.NewResource(
		vault.Name,
		resourceTypeVault,
		vault.ID,
		rs.WithAnnotation(&v2.ChildResourceType{ResourceTypeId: resourceTypeUser.Id}),
	)
	if err != nil {
		return nil, fmt.Errorf("baton-devolutions: failed to create vault resource: %w", err)
	}

	return ret, nil
}

func (v *vaultBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	bag := &pagination.Bag{}
	if err := bag.Unmarshal(pToken.Token); err != nil {
		return nil, "", nil, err
	}

	if bag.Current() == nil {
		bag.Push(pagination.PageState{
			ResourceTypeID: resourceTypeVault.Id,
			ResourceID:     "0",
		})
	}

	pageNumber, err := strconv.Atoi(bag.ResourceID())
	if err != nil {
		pageNumber = 0
	}

	resp, err := v.client.ListVaults(ctx, pageNumber, vaultsPageSize)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-devolutions: failed to list vaults: %w", err)
	}

	var resources []*v2.Resource
	for _, vault := range resp.Data {
		r, err := vaultResource(vault)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, r)
	}

	var nextPageToken string
	if resp.CurrentPage < resp.TotalPage {
		bag.Pop()
		bag.Push(pagination.PageState{
			ResourceTypeID: resourceTypeVault.Id,
			ResourceID:     strconv.Itoa(resp.CurrentPage + 1),
		})
		nextPageToken, err = bag.Marshal()
		if err != nil {
			return nil, "", nil, err
		}
	}

	return resources, nextPageToken, nil, nil
}

func (v *vaultBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var entitlements []*v2.Entitlement

	for _, perm := range vaultPermissions {
		permEnt := ent.NewPermissionEntitlement(
			resource,
			perm,
			ent.WithGrantableTo(resourceTypeUser),
			ent.WithDescription(fmt.Sprintf("%s access to %s vault", perm, resource.DisplayName)),
			ent.WithDisplayName(fmt.Sprintf("%s Vault %s", resource.DisplayName, perm)),
		)
		entitlements = append(entitlements, permEnt)
	}

	return entitlements, "", nil, nil
}

func (v *vaultBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	accessEntries, err := v.client.GetVaultAccess(ctx, resource.Id.Resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-devolutions: failed to get vault access: %w", err)
	}

	var grants []*v2.Grant
	for _, access := range accessEntries {
		if access.UserID == "" || access.PermissionSet == "" {
			continue
		}

		grant := sdkGrant.NewGrant(
			resource,
			access.PermissionSet,
			&v2.ResourceId{
				ResourceType: resourceTypeUser.Id,
				Resource:     access.UserID,
			},
		)
		grants = append(grants, grant)
	}

	return grants, "", nil, nil
}

func newVaultBuilder(client *client.Client) *vaultBuilder {
	return &vaultBuilder{
		resourceType: resourceTypeVault,
		client:       client,
	}
}
