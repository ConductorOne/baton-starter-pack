package main

import (
	"context"
	"fmt"
	"os"

	cfg "github.com/{{ cookiecutter.repo_owner }}/{{ cookiecutter.repo_name }}/pkg/config"
	"github.com/{{ cookiecutter.repo_owner }}/{{ cookiecutter.repo_name }}/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		"{{ cookiecutter.name }}",
		getConnector,
		cfg.Config,
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

func getConnector(ctx context.Context, cc *cfg.Connector) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	err := field.Validate(cfg.Config, cc)
	if err != nil {
		return nil, err
	}

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
