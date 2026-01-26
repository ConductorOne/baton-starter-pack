package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

// Configuration field definitions.
// Add your connector-specific fields here using field.StringField, field.BoolField, etc.
//
// Example:
//
//	var AccessToken = field.StringField(
//		"access-token",
//		field.WithDisplayName("Access Token"),
//		field.WithDescription("API access token for authentication"),
//		field.WithRequired(true),
//		field.WithIsSecret(true),
//	)

var (
	// ConfigurationFields defines the external configuration required for the
	// connector to run.
	ConfigurationFields = []field.SchemaField{
		// TODO: Add your fields here, e.g.:
		// AccessToken,
	}

	// FieldRelationships defines relationships between fields that can be
	// automatically validated. For example, a username and password can be
	// required together, or an access token can be marked as mutually exclusive
	// from the username/password pair.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run -tags=generate ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("{{ cookiecutter.display_name }}"),
)
