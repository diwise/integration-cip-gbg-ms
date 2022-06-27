package contextbroker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/diwise/context-broker/pkg/datamodels/fiware"
	"github.com/diwise/context-broker/pkg/ngsild/client"
	"github.com/diwise/context-broker/pkg/ngsild/types/entities"
	. "github.com/diwise/context-broker/pkg/ngsild/types/entities/decorators"
	"github.com/diwise/integration-cip-gbg-ms/internal/pkg/application/serviceguiden"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/logging"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/tracing"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type ContextBroker interface {
	QueryEntities(ctx context.Context, params url.Values) ([]byte, error)
	NewBeach(ctx context.Context, badplats serviceguiden.Content, nutsCode string) error
}

type contextBroker struct {
	baseUrl             string
	contextBrokerClient client.ContextBrokerClient
}

func New(log zerolog.Logger, brokerUrl string) ContextBroker {
	c := client.NewContextBrokerClient(brokerUrl)

	return &contextBroker{
		baseUrl:             brokerUrl,
		contextBrokerClient: c,
	}
}

var tracer = otel.Tracer("integration-cip-gbg-ms/context-broker")

func (c contextBroker) NewBeach(ctx context.Context, badplats serviceguiden.Content, nutsCode string) error {
	return newBeach(ctx, c.contextBrokerClient, badplats, nutsCode)
}

func newBeach(ctx context.Context, cbClient client.ContextBrokerClient, badplats serviceguiden.Content, nutsCode string) error {
	var id string

	if nutsCode != "" {
		id = fiware.BeachIDPrefix + nutsCode
	} else {
		id = fiware.BeachIDPrefix + badplats.Id
	}

	lat := badplats.Position.Latitude
	lon := badplats.Position.Longitude

	beach, err := fiware.NewBeach(
		id,
		badplats.Name,
		LocationMP([][][][]float64{{{
			{lon, lat},
			{lon, lat+0.0001},
			{lon+0.0001, lat+0.0001},
			{lon, lat},
		}}}),
		entities.DefaultContext(),
		Text("description", badplats.Description),
		Text("areaServed", badplats.AreaServed()),
		Text("dataProvider", "ServiceGuiden"),
		Text("source", badplats.Id),
		DateCreated(time.Now().UTC().Format(time.RFC3339)),
		TextList("facilities", badplats.Facilities()),
		TextList("seeAlso", textListWithoutEmptyValues([]string{badplats.SeeAlso(), getNutsCodeUrl(nutsCode), badplats.AccessibilityUrl})),
	)
	if err != nil {
		return err
	}

	headers := map[string][]string{"Content-Type": {"application/ld+json"}}
	_, err = cbClient.CreateEntity(ctx, beach, headers)

	return err
}

func updateBeach(ctx context.Context, cbClient client.ContextBrokerClient, badplats serviceguiden.Content, nutsCode string) error {
	return nil
}

func textListWithoutEmptyValues(values []string) []string {
	s := []string{}
	for _, v := range values {
		if v != "" {
			s = append(s, v)
		}
	}
	return s
}

func getNutsCodeUrl(nutsCode string) string {
	if nutsCode == "" {
		return ""
	}

	url := fmt.Sprintf("%s/testlocationprofile/%s", "https://badplatsen.havochvatten.se/badplatsen/api", nutsCode)
	return url
}

func (c contextBroker) QueryEntities(ctx context.Context, params url.Values) ([]byte, error) {
	var err error

	ctx, span := tracer.Start(ctx, "integration-cip-gbg/context-broker/queryentities")
	defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()

	log := logging.GetFromContext(ctx)

	httpClient := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	reqUrl := fmt.Sprintf("%s/%s?%s", c.baseUrl, "ngsi-ld/v1/entities", params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header = map[string][]string{
		"Accept": {"application/ld+json"},
		"Link":   {"<https://schema.lab.fiware.org/ld/context>; rel=\"http://www.w3.org/ns/json-ld#context\"; type=\"application/ld+json\""},
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve data from context-broker")
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("failed to retrieve data from context-broker, expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		return nil, fmt.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return nil, err
	}

	return body, err
}
