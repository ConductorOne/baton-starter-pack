package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	conn *Connector
}

func (u *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all users from the upstream service.
func (u *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	// TODO: Implement user listing with pagination
	// Example:
	//
	// page := ""
	// if pToken != nil && pToken.Token != "" {
	//     page = pToken.Token
	// }
	//
	// users, nextPage, err := u.conn.client.ListUsers(ctx, page, 100)
	// if err != nil {
	//     return nil, "", nil, fmt.Errorf("{{ cookiecutter.name }}: failed to list users: %w", err)
	// }
	//
	// var rv []*v2.Resource
	// for _, user := range users {
	//     resource, err := rs.NewUserResource(
	//         user.DisplayName,
	//         userResourceType,
	//         user.ID,
	//         []rs.UserTraitOption{
	//             rs.WithEmail(user.Email, true),
	//             rs.WithUserLogin(user.Username),
	//             rs.WithStatus(v2.UserTrait_Status_STATUS_ENABLED),
	//         },
	//         // ExternalId is CRITICAL for provisioning
	//         rs.WithExternalID(&v2.ExternalId{Id: user.ID}),
	//     )
	//     if err != nil {
	//         return nil, "", nil, err
	//     }
	//     rv = append(rv, resource)
	// }
	// return rv, nextPage, nil, nil

	_ = rs.NewUserResource
	_ = fmt.Sprintf
	return nil, "", nil, nil
}

// Entitlements returns an empty slice - users don't have entitlements.
func (u *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants returns an empty slice - users receive grants, they don't have them.
func (u *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// =============================================================================
// PROVISIONING: Create/Delete Users (ResourceManager interface)
// =============================================================================
// Uncomment and implement these methods to support user lifecycle management.
//
// func (u *userBuilder) Create(ctx context.Context, resource *v2.Resource) (*v2.Resource, annotations.Annotations, error) {
//     // TODO: Create user in upstream system
//     // Example:
//     //   userTrait, err := rs.GetUserTrait(resource)
//     //   if err != nil {
//     //       return nil, nil, err
//     //   }
//     //
//     //   newUser, err := u.conn.client.CreateUser(ctx, &CreateUserRequest{
//     //       Email:    userTrait.Email,
//     //       Username: userTrait.Login,
//     //   })
//     //   if err != nil {
//     //       return nil, nil, fmt.Errorf("{{ cookiecutter.name }}: failed to create user: %w", err)
//     //   }
//     //
//     //   return rs.NewUserResource(newUser.Name, userResourceType, newUser.ID, ...)
//     return nil, nil, fmt.Errorf("{{ cookiecutter.name }}: user creation not implemented")
// }
//
// func (u *userBuilder) Delete(ctx context.Context, resourceId *v2.ResourceId) (annotations.Annotations, error) {
//     // TODO: Delete user from upstream system
//     // Example:
//     //   err := u.conn.client.DeleteUser(ctx, resourceId.Resource)
//     //   if err != nil {
//     //       return nil, fmt.Errorf("{{ cookiecutter.name }}: failed to delete user: %w", err)
//     //   }
//     //   return nil, nil
//     return nil, fmt.Errorf("{{ cookiecutter.name }}: user deletion not implemented")
// }

func newUserBuilder(conn *Connector) *userBuilder {
	return &userBuilder{conn: conn}
}
