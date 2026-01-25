package main

import (
	"context"
	"fmt"
	"os"

	"github.com/{{ cookiecutter.repo_owner }}/{{ cookiecutter.repo_name }}/pkg/connector"
	configSdk "github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var version = "dev"

// Config holds the connector configuration.
type Config struct {
	// Add connector-specific fields here.
	// Example: APIKey string `mapstructure:"api-key"`
}

// Implement field.Configurable interface.
func (c *Config) GetString(key string) string       { return "" }
func (c *Config) GetBool(key string) bool           { return false }
func (c *Config) GetInt(key string) int             { return 0 }
func (c *Config) GetStringSlice(key string) []string { return nil }
func (c *Config) GetStringMap(key string) map[string]any { return nil }

// Configuration fields for the connector.
var configFields = []field.SchemaField{
	// TODO: Add your connector-specific fields here, e.g.:
	// field.StringField("api-key", field.WithRequired(true), field.WithDescription("API key for authentication")),
}

// ConfigSchema is the configuration schema for the connector.
var ConfigSchema = field.NewConfiguration(configFields)

func main() {
	ctx := context.Background()

	_, cmd, err := configSdk.DefineConfiguration(
		ctx,
		"{{ cookiecutter.name }}",
		getConnector,
		ConfigSchema,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, cfg *Config) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	cb, err := connector.New(ctx)
	if err != nil {
		l.Error("{{ cookiecutter.name }}: error creating connector", zap.Error(err))
		return nil, fmt.Errorf("{{ cookiecutter.name }}: failed to create connector: %w", err)
	}

	c, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("{{ cookiecutter.name }}: error wrapping connector", zap.Error(err))
		return nil, fmt.Errorf("{{ cookiecutter.name }}: failed to initialize connector: %w", err)
	}

	return c, nil
}
