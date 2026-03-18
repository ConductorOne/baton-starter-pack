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
	groupsPageSize    = 50
	memberEntitlement = "member"
)

type groupBuilder struct {
	resourceType *v2.ResourceType
	client       *client.Client
}

func (g *groupBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return g.resourceType
}

func groupResource(group client.UserGroup) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"group_name": group.Name,
		"group_id":   group.ID,
	}

	if group.Description != "" {
		profile["description"] = group.Description
	}

	ret, err := rs.NewGroupResource(
		group.Name,
		resourceTypeGroup,
		group.ID,
		[]rs.GroupTraitOption{
			rs.WithGroupProfile(profile),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("baton-devolutions: failed to create group resource: %w", err)
	}

	return ret, nil
}

func (g *groupBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	bag := &pagination.Bag{}
	if err := bag.Unmarshal(pToken.Token); err != nil {
		return nil, "", nil, err
	}

	if bag.Current() == nil {
		bag.Push(pagination.PageState{
			ResourceTypeID: resourceTypeGroup.Id,
			ResourceID:     "0",
		})
	}

	pageNumber, err := strconv.Atoi(bag.ResourceID())
	if err != nil {
		pageNumber = 0
	}

	resp, err := g.client.ListGroups(ctx, pageNumber, groupsPageSize)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-devolutions: failed to list groups: %w", err)
	}

	var resources []*v2.Resource
	for _, group := range resp.Data {
		r, err := groupResource(group)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, r)
	}

	var nextPageToken string
	if resp.CurrentPage < resp.TotalPage {
		bag.Pop()
		bag.Push(pagination.PageState{
			ResourceTypeID: resourceTypeGroup.Id,
			ResourceID:     strconv.Itoa(resp.CurrentPage + 1),
		})
		nextPageToken, err = bag.Marshal()
		if err != nil {
			return nil, "", nil, err
		}
	}

	return resources, nextPageToken, nil, nil
}

func (g *groupBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	memberEnt := ent.NewAssignmentEntitlement(
		resource,
		memberEntitlement,
		ent.WithGrantableTo(resourceTypeUser),
		ent.WithDescription(fmt.Sprintf("Member of %s group", resource.DisplayName)),
		ent.WithDisplayName(fmt.Sprintf("%s Group Member", resource.DisplayName)),
	)

	return []*v2.Entitlement{memberEnt}, "", nil, nil
}

func (g *groupBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	members, err := g.client.GetGroupMembers(ctx, resource.Id.Resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-devolutions: failed to get group members: %w", err)
	}

	var grants []*v2.Grant
	for _, member := range members {
		grant := sdkGrant.NewGrant(
			resource,
			memberEntitlement,
			&v2.ResourceId{
				ResourceType: resourceTypeUser.Id,
				Resource:     member.UserID,
			},
		)
		grants = append(grants, grant)
	}

	return grants, "", nil, nil
}

func newGroupBuilder(client *client.Client) *groupBuilder {
	return &groupBuilder{
		resourceType: resourceTypeGroup,
		client:       client,
	}
}
