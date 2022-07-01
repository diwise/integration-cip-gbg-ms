package main

import (
	"context"
	"flag"

	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application"
	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/contextbroker"
	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/lookup"
	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/serviceguiden"
	"github.com/diwise/service-chassis/pkg/infrastructure/buildinfo"
	"github.com/diwise/service-chassis/pkg/infrastructure/env"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y"
)

var lookupTableFilePath string
var serviceGuidenFilePath string

const serviceName string = "integration-cip-gbg"

func main() {
	serviceVersion := buildinfo.SourceVersion()

	ctx, logger, cleanup := o11y.Init(context.Background(), serviceName, serviceVersion)
	defer cleanup()

	flag.StringVar(&lookupTableFilePath, "references", "/opt/diwise/config/lookup.csv", "A file with cross-references from service guiden to nutscodes and devices")
	flag.StringVar(&serviceGuidenFilePath, "sg", "/opt/diwise/config/serviceguiden.json", "A file with ServiceGuiden contents")
	flag.Parse()

	serviceGuidenUrl := env.GetVariableOrDefault(logger, "SERVICE_GUIDEN", "https://microservices.goteborg.se/sdw-service/api/internal/v1/sites?size=10000")
	contextBrokerUrl := env.GetVariableOrDefault(logger, "CONTEXT_BROKER", "http://context-broker")

	sgClient := serviceguiden.New(serviceGuidenUrl, serviceGuidenFilePath)
	cbClient := contextbroker.New(logger, contextBrokerUrl)
	lookupTable := lookup.New(logger, lookupTableFilePath)

	app := application.New(sgClient, cbClient, lookupTable, logger)

	err := app.CreateOrUpdateBeachModels(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("beach update failed")
	}
}
