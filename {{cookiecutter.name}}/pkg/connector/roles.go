package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const assignedEntitlement = "assigned"

type roleBuilder struct {
	conn *Connector
}

func (r *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return roleResourceType
}

// List returns all roles from the upstream service.
func (r *roleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	// TODO: Implement role listing
	_ = rs.NewRoleResource
	_ = fmt.Sprintf
	return nil, "", nil, nil
}

// Entitlements returns the "assigned" entitlement for the role.
func (r *roleBuilder) Entitlements(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	entitlement := ent.NewAssignmentEntitlement(
		resource,
		assignedEntitlement,
		ent.WithGrantableTo(userResourceType),
		ent.WithDisplayName(fmt.Sprintf("%s Role", resource.DisplayName)),
		ent.WithDescription(fmt.Sprintf("Assigned the %s role", resource.DisplayName)),
	)
	return []*v2.Entitlement{entitlement}, "", nil, nil
}

// Grants returns all users who are assigned this role.
func (r *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	// TODO: Implement role assignment listing
	_ = grant.NewGrant
	return nil, "", nil, nil
}

// =============================================================================
// PROVISIONING: Grant/Revoke role assignment (ResourceProvisioner interface)
// =============================================================================
// Uncomment and implement these methods to support role assignment management.
//
// func (r *roleBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
//     // TODO: Assign role to user
//     return nil, fmt.Errorf("{{ cookiecutter.name }}: grant not implemented")
// }
//
// func (r *roleBuilder) Revoke(ctx context.Context, grantToRevoke *v2.Grant) (annotations.Annotations, error) {
//     // TODO: Unassign role from user
//     return nil, fmt.Errorf("{{ cookiecutter.name }}: revoke not implemented")
// }

func newRoleBuilder(conn *Connector) *roleBuilder {
	return &roleBuilder{conn: conn}
}
