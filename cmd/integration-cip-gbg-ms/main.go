package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"log/slog"

	"github.com/diwise/context-broker/pkg/datamodels/fiware"
	"github.com/diwise/context-broker/pkg/ngsild/client"
	"github.com/google/uuid"

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

	logger.Debug("args:", slog.String("references", lookupTableFilePath), slog.String("sg", serviceGuidenFilePath))

	serviceGuidenUrl := env.GetVariableOrDefault(ctx, "SERVICE_GUIDEN", "https://microservices.goteborg.se/sdw-service/api/internal/v1/sites?size=10000")
	contextBrokerUrl := env.GetVariableOrDefault(ctx, "CONTEXT_BROKER", "http://context-broker")

	logger.Debug("env:", slog.String("SERVICE_GUIDEN", serviceGuidenUrl), slog.String("CONTEXT_BROKER", contextBrokerUrl))

	cbClient := client.NewContextBrokerClient(contextBrokerUrl)
	sgClient := serviceguiden.New(ctx, serviceGuidenUrl, serviceGuidenFilePath)
	lookupTable := lookup.New(logger, lookupTableFilePath)

	err := run(ctx, sgClient, lookupTable, cbClient, logger)
	if err != nil {
		logger.Error("failed to create or update beaches", "err", err.Error())
	}
}

func run(ctx context.Context, sgClient serviceguiden.ServiceGuidenClient, lookupTable lookup.LookupTable, cbClient client.ContextBrokerClient, logger *slog.Logger) error {
	badplatser, err := sgClient.Badplatser(ctx)
	if err != nil {
		return err
	}

	errs := []error{}

	for _, badplats := range badplatser {
		nutsCode, _ := lookupTable.GetNutsCode(badplats.ID())
		props := cip.NewBeachProps(badplats, nutsCode)
		beachID := fiware.BeachIDPrefix + deterministicGUID("ServiceGuiden", badplats.ID())

		err := cip.MergeOrCreate(ctx, cbClient, beachID, fiware.BeachTypeName, props)
		if err != nil {
			logger.Error("faild to merge beach", slog.String("beach_id", beachID), slog.String("err", err.Error()))
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func deterministicGUID(dataProvider string, id string) string {
	md5hash := md5.New()
	md5hash.Write([]byte(id + dataProvider))
	md5string := hex.EncodeToString(md5hash.Sum(nil))

	unique, err := uuid.FromBytes([]byte(md5string[0:16]))
	if err != nil {
		return uuid.New().String()
	}

	return unique.String()
}
