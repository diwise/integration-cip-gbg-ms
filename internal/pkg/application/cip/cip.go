package cip

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/diwise/context-broker/pkg/ngsild/client"
	ngsierrors "github.com/diwise/context-broker/pkg/ngsild/errors"
	"github.com/diwise/context-broker/pkg/ngsild/types/entities"
	"github.com/diwise/context-broker/pkg/ngsild/types/entities/decorators"
	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/serviceguiden"
	"github.com/diwise/service-chassis/pkg/infrastructure/env"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/logging"
	"github.com/rs/zerolog"
)

var havOchVattenProfileUrl string
var seeAlsoUrl string
var dataProvider string
var source string

func init() {
	havOchVattenProfileUrl = env.GetVariableOrDefault(zerolog.Logger{}, "HAV_OCH_VATTEN_PROFILE_URL", "https://badplatsen.havochvatten.se/badplatsen/api/testlocationprofile")
	seeAlsoUrl = env.GetVariableOrDefault(zerolog.Logger{}, "SEE_ALSO_URL", "https://goteborg.se/wps/portal/start/uppleva-och-gora/idrott-motion-och-friluftsliv/simma-och-bada/badplatser/hitta-badplatser-utomhusbad/?id=")
	dataProvider = env.GetVariableOrDefault(zerolog.Logger{}, "DATA_PROVIDER", "ServiceGuiden")
	source = env.GetVariableOrDefault(zerolog.Logger{}, "SOURCE", "se:goteborg:serviceguiden:businessid:")
}

func MergeOrCreate(ctx context.Context, cbClient client.ContextBrokerClient, id string, typeName string, properties []entities.EntityDecoratorFunc) error {
	log := logging.GetFromContext(ctx)

	headers := map[string][]string{"Content-Type": {"application/ld+json"}}

	fragment, err := entities.NewFragment(properties...)
	if err != nil {
		return fmt.Errorf("failed to create new fragment for entity %s, %w", id, err)
	}

	_, err = cbClient.MergeEntity(ctx, id, fragment, headers)
	if err != nil {
		if !errors.Is(err, ngsierrors.ErrNotFound) {
			return fmt.Errorf("failed to merge entity %s, %w", id, err)
		}

		properties = append(properties, entities.DefaultContext())

		entity, err := entities.New(id, typeName, properties...)
		if err != nil {
			return fmt.Errorf("failed to create new entity props for entity %s, %w", id, err)
		}

		_, err = cbClient.CreateEntity(ctx, entity, headers)
		if err != nil {
			return fmt.Errorf("failed to create entity %s, %w", id, err)
		}

		log.Debug().Msgf("create entity %s", id)

		return nil
	}

	log.Debug().Msgf("merge entity %s", id)

	return nil
}

func NewBeachProps(badplats serviceguiden.Content, nutsCode string) []entities.EntityDecoratorFunc {
	props := []entities.EntityDecoratorFunc{}

	lat := badplats.Position.Latitude
	lon := badplats.Position.Longitude

	seeAlso := filter([]string{getSeeAlso(badplats), getNutsCodeUrl(nutsCode), badplats.AccessibilityUrl}, func(s string) bool {
		return s != ""
	})

	source := fmt.Sprintf("%s%d", source, badplats.BusinessId)

	props = append(props,
		decorators.LocationMP([][][][]float64{{{
			{lon, lat},
			{lon, lat + 0.0001},
			{lon + 0.0001, lat + 0.0001},
			{lon, lat},
		}}}),
		entities.DefaultContext(),
		decorators.Name(badplats.Name),
		decorators.Text("description", badplats.Description),
		decorators.Text("areaServed", badplats.AreaServed()),
		decorators.Text("dataProvider", dataProvider),
		decorators.Text("source", source),
		decorators.DateCreated(time.Now().UTC().Format(time.RFC3339)),
		decorators.TextList("beachType", badplats.BeachTypes()),
		decorators.TextList("seeAlso", seeAlso),
	)

	return props
}

func getSeeAlso(badplats serviceguiden.Content) string {
	return fmt.Sprintf("%s%d", seeAlsoUrl, badplats.BusinessId)
}

func getNutsCodeUrl(nutsCode string) string {
	if nutsCode == "" {
		return ""
	}

	url := fmt.Sprintf("%s/%s", havOchVattenProfileUrl, nutsCode)
	return url
}

func filter[T any](data []T, f func(T) bool) []T {
	fltd := make([]T, 0)

	for _, e := range data {
		if f(e) {
			fltd = append(fltd, e)
		}
	}

	return fltd
}
