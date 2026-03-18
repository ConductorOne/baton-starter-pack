package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-devolutions/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const assignedEntitlement = "assigned"

type roleBuilder struct {
	resourceType *v2.ResourceType
	client       *client.Client
}

func (r *roleBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return r.resourceType
}

func roleResource(role client.Role) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"role_name": role.Name,
		"role_id":   role.ID,
	}

	if role.Description != "" {
		profile["description"] = role.Description
	}

	ret, err := rs.NewRoleResource(
		role.Name,
		resourceTypeRole,
		role.ID,
		[]rs.RoleTraitOption{
			rs.WithRoleProfile(profile),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("baton-devolutions: failed to create role resource: %w", err)
	}

	return ret, nil
}

func (r *roleBuilder) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	roles, err := r.client.ListRoles(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-devolutions: failed to list roles: %w", err)
	}

	var resources []*v2.Resource
	for _, role := range roles {
		res, err := roleResource(role)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, res)
	}

	return resources, "", nil, nil
}

func (r *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	assignedEnt := ent.NewAssignmentEntitlement(
		resource,
		assignedEntitlement,
		ent.WithGrantableTo(resourceTypeUser),
		ent.WithDescription(fmt.Sprintf("Assigned the %s role", resource.DisplayName)),
		ent.WithDisplayName(fmt.Sprintf("%s Role Assigned", resource.DisplayName)),
	)

	return []*v2.Entitlement{assignedEnt}, "", nil, nil
}

func (r *roleBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	// Role grants are expressed via vault access permission sets.
	// Direct role-to-user mapping is not available through the DVLS API;
	// instead, roles manifest as permission sets on vault access entries.
	return nil, "", nil, nil
}

func newRoleBuilder(client *client.Client) *roleBuilder {
	return &roleBuilder{
		resourceType: resourceTypeRole,
		client:       client,
	}
}
