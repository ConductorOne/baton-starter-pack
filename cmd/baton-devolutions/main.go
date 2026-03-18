package main

import (
	"context"
	"fmt"
	"os"

	devConfig "github.com/conductorone/baton-devolutions/pkg/config"
	"github.com/conductorone/baton-devolutions/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	connectorName = "baton-devolutions"
	version       = "dev"
)

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		connectorName,
		getConnector,
		devConfig.ConfigurationSchema,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, v *viper.Viper) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	serverURL := v.GetString(devConfig.ServerURLField.FieldName)
	appKey := v.GetString(devConfig.AppKeyField.FieldName)
	appSecret := v.GetString(devConfig.AppSecretField.FieldName)

	cb, err := connector.New(ctx, serverURL, appKey, appSecret)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	c, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector server", zap.Error(err))
		return nil, err
	}

	return c, nil
}
