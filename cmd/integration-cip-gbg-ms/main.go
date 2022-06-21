package main

import (
	"context"
	"flag"

	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application"
	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/contextbroker"
	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/serviceguiden"
	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/lookup"
	"github.com/diwise/service-chassis/pkg/infrastructure/buildinfo"
	"github.com/diwise/service-chassis/pkg/infrastructure/env"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y"
)

var lookupTableFilePath string

func main() {
	serviceVersion := buildinfo.SourceVersion()
	serviceName := "integration-cip-gbg"

	ctx, logger, cleanup := o11y.Init(context.Background(), serviceName, serviceVersion)
	defer cleanup()

	flag.StringVar(&lookupTableFilePath, "references", "/opt/diwise/config/lookup.csv", "A file with cross-references from service guiden to nutscodes and devices")
	flag.Parse()

	serviceGuidenUrl := env.GetVariableOrDefault(logger, "SERVICE_GUIDEN", "https://microservices.goteborg.se/sdw-service/api/internal/v1/sites?size=10000")
	contextBrokerUrl := env.GetVariableOrDefault(logger, "CONTEXT_BROKER", "http://context-broker:8080")

	sgClient := serviceguiden.New(serviceGuidenUrl)
	cbClient := contextbroker.New(logger, contextBrokerUrl)
	lookupTable := lookup.New(logger, lookupTableFilePath)

	app := application.New(sgClient, cbClient, lookupTable)

	app.CreateOrUpdateBeachModels(ctx)
}
