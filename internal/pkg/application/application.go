package application

import (
	"context"

	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/contextbroker"
	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/lookup"
	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/serviceguiden"
)

type App interface {
	CreateOrUpdateBeachModels(ctx context.Context) error
}

type app struct {
	sgClient serviceguiden.ServiceGuidenClient
	cbClient contextbroker.ContextBroker
	lookup lookup.LookupTable
}

func New(sgClient serviceguiden.ServiceGuidenClient, cbClient contextbroker.ContextBroker, lookup lookup.LookupTable) App {
	return &app{
		sgClient: sgClient,
		cbClient: cbClient,
		lookup: lookup,
	}
}

func (a app) CreateOrUpdateBeachModels(ctx context.Context) error {
	badplatser, err := a.sgClient.Badplatser(ctx)
	if err != nil {
		return err
	}

	for _, badplats := range badplatser {
		nutsCode, _ := a.lookup.GetNutsCode(badplats.Id)
		a.cbClient.NewBeach(ctx, badplats, nutsCode)
	}

	return nil
}
