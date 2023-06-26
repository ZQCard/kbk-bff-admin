//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/ZQCard/kbk-bff-admin/internal/conf"
	"github.com/ZQCard/kbk-bff-admin/internal/data"
	"github.com/ZQCard/kbk-bff-admin/internal/server"
	"github.com/ZQCard/kbk-bff-admin/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// wireApp init kratos application.
func wireApp(*conf.Env, *conf.Server, *conf.Service, *conf.Data, *conf.Registry, *conf.Bootstrap, *conf.Auth, *conf.Endpoint, log.Logger, *tracesdk.TracerProvider) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, service.ProviderSet, newApp))
}
