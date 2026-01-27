package connector

import (
	"context"
	"fmt"

	config "github.com/conductorone/baton-sdk/pb/c1/config/v1"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/actions"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

// =============================================================================
// BATON ACTIONS: Custom operations exposed to ConductorOne
// =============================================================================
//
// Actions are arbitrary operations your connector can perform. Unlike Grant/Revoke
// (which modify access) or Create/Delete (which manage resources), Actions are
// general-purpose operations that ConductorOne can trigger.
//
// Common action types:
// - ACTION_TYPE_ACCOUNT_ENABLE / ACTION_TYPE_ACCOUNT_DISABLE
// - ACTION_TYPE_ACCOUNT_UPDATE_PROFILE
// - Custom operations specific to your system
//
// To enable actions, implement the GlobalActions method below.

// Example action schema for disabling an account
var disableAccountAction = &v2.BatonActionSchema{
	Name: "disableAccount",
	Arguments: []*config.Field{
		{
			Name:        "accountId",
			DisplayName: "Account ID",
			Description: "The ID of the account to disable",
			Field:       &config.Field_StringField{},
			IsRequired:  true,
		},
	},
	ReturnTypes: []*config.Field{
		{
			Name:        "success",
			DisplayName: "Success",
			Field:       &config.Field_BoolField{},
		},
	},
	ActionType: []v2.ActionType{
		v2.ActionType_ACTION_TYPE_ACCOUNT,
		v2.ActionType_ACTION_TYPE_ACCOUNT_DISABLE,
	},
}

// Example action schema for enabling an account
var enableAccountAction = &v2.BatonActionSchema{
	Name: "enableAccount",
	Arguments: []*config.Field{
		{
			Name:        "accountId",
			DisplayName: "Account ID",
			Description: "The ID of the account to enable",
			Field:       &config.Field_StringField{},
			IsRequired:  true,
		},
	},
	ReturnTypes: []*config.Field{
		{
			Name:        "success",
			DisplayName: "Success",
			Field:       &config.Field_BoolField{},
		},
	},
	ActionType: []v2.ActionType{
		v2.ActionType_ACTION_TYPE_ACCOUNT,
		v2.ActionType_ACTION_TYPE_ACCOUNT_ENABLE,
	},
}

// GlobalActions registers custom actions with the SDK.
// Uncomment and implement to enable actions.
//
// func (c *Connector) GlobalActions(ctx context.Context, registry actions.ActionRegistry) error {
//     if err := registry.Register(ctx, disableAccountAction, c.disableAccount); err != nil {
//         return err
//     }
//     if err := registry.Register(ctx, enableAccountAction, c.enableAccount); err != nil {
//         return err
//     }
//     return nil
// }

func (c *Connector) disableAccount(ctx context.Context, args *structpb.Struct) (*structpb.Struct, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	accountId, ok := args.Fields["accountId"]
	if !ok {
		return nil, nil, fmt.Errorf("missing required argument accountId")
	}

	l.Info("{{ cookiecutter.name }}: disabling account", zap.String("accountId", accountId.GetStringValue()))

	// TODO: Implement account disabling
	// Example:
	//   err := c.client.DisableUser(ctx, accountId.GetStringValue())
	//   if err != nil {
	//       return nil, nil, fmt.Errorf("{{ cookiecutter.name }}: failed to disable account: %w", err)
	//   }

	response := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"success": structpb.NewBoolValue(true),
		},
	}
	return response, nil, nil
}

func (c *Connector) enableAccount(ctx context.Context, args *structpb.Struct) (*structpb.Struct, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	accountId, ok := args.Fields["accountId"]
	if !ok {
		return nil, nil, fmt.Errorf("missing required argument accountId")
	}

	l.Info("{{ cookiecutter.name }}: enabling account", zap.String("accountId", accountId.GetStringValue()))

	// TODO: Implement account enabling
	// Example:
	//   err := c.client.EnableUser(ctx, accountId.GetStringValue())
	//   if err != nil {
	//       return nil, nil, fmt.Errorf("{{ cookiecutter.name }}: failed to enable account: %w", err)
	//   }

	response := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"success": structpb.NewBoolValue(true),
		},
	}
	return response, nil, nil
}

// Ensure imports are used (remove after implementing)
var _ = actions.ActionRegistry(nil)
