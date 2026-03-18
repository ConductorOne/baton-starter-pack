package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-devolutions/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const usersPageSize = 50

type userBuilder struct {
	resourceType *v2.ResourceType
	client       *client.Client
}

func (u *userBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return u.resourceType
}

func userResource(user client.User) (*v2.Resource, error) {
	var userStatus v2.UserTrait_Status_Status
	if user.IsEnabled {
		userStatus = v2.UserTrait_Status_STATUS_ENABLED
	} else {
		userStatus = v2.UserTrait_Status_STATUS_DISABLED
	}

	profile := map[string]interface{}{
		"first_name":          user.FirstName,
		"last_name":           user.LastName,
		"login":               user.Username,
		"user_id":             user.ID,
		"user_type":           user.UserType,
		"authentication_type": user.AuthenticationType,
		"is_administrator":    user.IsAdministrator,
	}

	userTraitOptions := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(userStatus),
		rs.WithUserLogin(user.Username),
	}

	if user.Email != "" {
		userTraitOptions = append(userTraitOptions, rs.WithEmail(user.Email, true))
	}

	ret, err := rs.NewUserResource(
		displayName(user),
		resourceTypeUser,
		user.ID,
		userTraitOptions,
	)
	if err != nil {
		return nil, fmt.Errorf("baton-devolutions: failed to create user resource: %w", err)
	}

	return ret, nil
}

func displayName(user client.User) string {
	if user.FirstName != "" || user.LastName != "" {
		name := user.FirstName
		if user.LastName != "" {
			if name != "" {
				name += " "
			}
			name += user.LastName
		}
		return name
	}
	return user.Username
}

func (u *userBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	bag := &pagination.Bag{}
	if err := bag.Unmarshal(pToken.Token); err != nil {
		return nil, "", nil, err
	}

	if bag.Current() == nil {
		bag.Push(pagination.PageState{
			ResourceTypeID: resourceTypeUser.Id,
			ResourceID:     "0",
		})
	}

	pageNumber, err := strconv.Atoi(bag.ResourceID())
	if err != nil {
		pageNumber = 0
	}

	resp, err := u.client.ListUsers(ctx, pageNumber, usersPageSize)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-devolutions: failed to list users: %w", err)
	}

	var resources []*v2.Resource
	for _, user := range resp.Data {
		r, err := userResource(user)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, r)
	}

	var nextPageToken string
	if resp.CurrentPage < resp.TotalPage {
		bag.Pop()
		bag.Push(pagination.PageState{
			ResourceTypeID: resourceTypeUser.Id,
			ResourceID:     strconv.Itoa(resp.CurrentPage + 1),
		})
		nextPageToken, err = bag.Marshal()
		if err != nil {
			return nil, "", nil, err
		}
	}

	return resources, nextPageToken, nil, nil
}

func (u *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (u *userBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(client *client.Client) *userBuilder {
	return &userBuilder{
		resourceType: resourceTypeUser,
		client:       client,
	}
}
