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

const memberEntitlement = "member"

type groupBuilder struct {
	conn *Connector
}

func (g *groupBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return groupResourceType
}

// List returns all groups from the upstream service.
func (g *groupBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	// TODO: Implement group listing
	_ = rs.NewGroupResource
	_ = fmt.Sprintf
	return nil, "", nil, nil
}

// Entitlements returns the "member" entitlement for the group.
func (g *groupBuilder) Entitlements(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	entitlement := ent.NewAssignmentEntitlement(
		resource,
		memberEntitlement,
		ent.WithGrantableTo(userResourceType),
		ent.WithDisplayName(fmt.Sprintf("%s Member", resource.DisplayName)),
		ent.WithDescription(fmt.Sprintf("Member of the %s group", resource.DisplayName)),
	)
	return []*v2.Entitlement{entitlement}, "", nil, nil
}

// Grants returns all users who are members of this group.
func (g *groupBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	// TODO: Implement group membership listing
	// Example:
	//   members, nextPage, err := g.conn.client.ListGroupMembers(ctx, resource.Id.Resource, ...)
	//   for _, member := range members {
	//       rv = append(rv, grant.NewGrant(resource, memberEntitlement, &v2.ResourceId{
	//           ResourceType: userResourceType.Id,
	//           Resource:     member.UserID,
	//       }))
	//   }
	_ = grant.NewGrant
	return nil, "", nil, nil
}

// =============================================================================
// PROVISIONING: Grant/Revoke group membership (ResourceProvisioner interface)
// =============================================================================
// Uncomment and implement these methods to support group membership management.
//
// func (g *groupBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
//     // TODO: Add user to group
//     // Example:
//     //   groupID := entitlement.Resource.Id.Resource
//     //   userID := principal.Id.Resource
//     //   err := g.conn.client.AddGroupMember(ctx, groupID, userID)
//     //   if err != nil {
//     //       return nil, fmt.Errorf("{{ cookiecutter.name }}: failed to grant membership: %w", err)
//     //   }
//     //   return nil, nil
//     return nil, fmt.Errorf("{{ cookiecutter.name }}: grant not implemented")
// }
//
// func (g *groupBuilder) Revoke(ctx context.Context, grantToRevoke *v2.Grant) (annotations.Annotations, error) {
//     // TODO: Remove user from group
//     // Example:
//     //   groupID := grantToRevoke.Entitlement.Resource.Id.Resource
//     //   userID := grantToRevoke.Principal.Id.Resource
//     //   err := g.conn.client.RemoveGroupMember(ctx, groupID, userID)
//     //   if err != nil {
//     //       return nil, fmt.Errorf("{{ cookiecutter.name }}: failed to revoke membership: %w", err)
//     //   }
//     //   return nil, nil
//     return nil, fmt.Errorf("{{ cookiecutter.name }}: revoke not implemented")
// }

func newGroupBuilder(conn *Connector) *groupBuilder {
	return &groupBuilder{conn: conn}
}
