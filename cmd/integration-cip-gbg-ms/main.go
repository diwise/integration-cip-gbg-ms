package main

import (
	"context"
	"flag"

	"github.com/diwise/context-broker/pkg/datamodels/fiware"
	"github.com/diwise/context-broker/pkg/ngsild/client"
	"github.com/rs/zerolog"

	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/cip"
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

	cbClient := client.NewContextBrokerClient(contextBrokerUrl)
	sgClient := serviceguiden.New(serviceGuidenUrl, serviceGuidenFilePath)
	lookupTable := lookup.New(logger, lookupTableFilePath)

	err := run(ctx, sgClient, lookupTable, cbClient, logger)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create or update beaches")
	}
}

func run(ctx context.Context, sgClient serviceguiden.ServiceGuidenClient, lookupTable lookup.LookupTable, cbClient client.ContextBrokerClient, logger zerolog.Logger) error {
	badplatser, err := sgClient.Badplatser(ctx)
	if err != nil {
		return err
	}

	getBeachID := func(nutsCode, badplatsID string) string {
		if nutsCode != "" {
			return fiware.BeachIDPrefix + nutsCode
		} else {
			return fiware.BeachIDPrefix + badplatsID
		}
	}

	for _, badplats := range badplatser {
		nutsCode, _ := lookupTable.GetNutsCode(badplats.Id)
		props := cip.NewBeach(badplats, nutsCode)
		beachID := getBeachID(nutsCode, badplats.Id)

		err := cip.MergeOrCreate(ctx, cbClient, beachID, fiware.BeachTypeName, props)
		if err != nil {
			logger.Error().Err(err).Msgf("faild to merge %s", beachID)
			return err
		}
	}

	return nil
}
