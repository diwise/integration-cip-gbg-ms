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
	"github.com/rs/zerolog"
)

var havOchVattenProfileUrl string
var seeAlsoUrl string

func init() {
	havOchVattenProfileUrl = env.GetVariableOrDefault(zerolog.Logger{}, "HAV_OCH_VATTEN_PROFILE_URL", "https://badplatsen.havochvatten.se/badplatsen/api/testlocationprofile")
	seeAlsoUrl = env.GetVariableOrDefault(zerolog.Logger{}, "SEE_ALSO_URL", "https://goteborg.se/wps/portal/start/kultur-och-fritid/fritid-och-natur/friluftsliv-natur-och/badplatser--utomhusbad/badplatser-utomhusbad/?id=")
}

func MergeOrCreate(ctx context.Context, cbClient client.ContextBrokerClient, id string, typeName string, properties []entities.EntityDecoratorFunc) error {
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
	}

	return nil
}

func NewBeach(badplats serviceguiden.Content, nutsCode string) []entities.EntityDecoratorFunc {
	props := []entities.EntityDecoratorFunc{}

	lat := badplats.Position.Latitude
	lon := badplats.Position.Longitude

	seeAlso := getSeeAlso(badplats)

	props = append(props,
		decorators.LocationMP([][][][]float64{{{
			{lon, lat},
			{lon, lat + 0.0001},
			{lon + 0.0001, lat + 0.0001},
			{lon, lat},
		}}}),
		entities.DefaultContext(),
		decorators.Text("description", badplats.Description),
		decorators.Text("areaServed", badplats.AreaServed()),
		decorators.Text("dataProvider", "ServiceGuiden"),
		decorators.Text("source", badplats.Id),
		decorators.DateCreated(time.Now().UTC().Format(time.RFC3339)),
		decorators.TextList("beachType", badplats.BeachTypes()),
		decorators.TextList("seeAlso", filter([]string{seeAlso, getNutsCodeUrl(nutsCode), badplats.AccessibilityUrl}, func(s string) bool {
			return s != ""
		})),
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
