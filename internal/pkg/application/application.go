package application

import (
	"context"

	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/contextbroker"
	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/lookup"
	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/serviceguiden"
	"github.com/rs/zerolog"
)

type App interface {
	CreateOrUpdateBeachModels(ctx context.Context) error
}

type app struct {
	sgClient serviceguiden.ServiceGuidenClient
	cbClient contextbroker.ContextBroker
	lookup   lookup.LookupTable
	log      zerolog.Logger
}

func New(sgClient serviceguiden.ServiceGuidenClient, cbClient contextbroker.ContextBroker, lookup lookup.LookupTable, logger zerolog.Logger) App {
	return &app{
		sgClient: sgClient,
		cbClient: cbClient,
		lookup:   lookup,
		log:      logger,
	}
}

func (a app) CreateOrUpdateBeachModels(ctx context.Context) error {
	badplatser, err := a.sgClient.Badplatser(ctx)
	if err != nil {
		return err
	}

	for _, badplats := range badplatser {
		nutsCode, _ := a.lookup.GetNutsCode(badplats.Id)
		if err = a.cbClient.NewBeach(ctx, badplats, nutsCode); err != nil {
			a.log.Err(err).Msg("unable to create new Beach")
		}
		/*
			if err == "allready exists" {
				a.cbClient.UpdateBeach(ctx, badplats, nutsCode)
			}
		*/
	}

	return nil
}
