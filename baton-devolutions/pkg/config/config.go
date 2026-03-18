package config

import "github.com/conductorone/baton-sdk/pkg/field"

var (
	ServerURLField = field.StringField(
		"server-url",
		field.WithDescription("Devolutions Server URL (e.g., https://dvls.example.com)"),
		field.WithRequired(true),
	)

	AppKeyField = field.StringField(
		"app-key",
		field.WithDescription("Application Identity key for DVLS authentication"),
		field.WithRequired(true),
	)

	AppSecretField = field.StringField(
		"app-secret",
		field.WithDescription("Application Identity secret for DVLS authentication"),
		field.WithRequired(true),
		field.WithIsSecret(true),
	)

	ConfigurationFields = []field.SchemaField{
		ServerURLField,
		AppKeyField,
		AppSecretField,
	}

	ConfigurationSchema = field.Configuration{
		Fields: ConfigurationFields,
	}
)
