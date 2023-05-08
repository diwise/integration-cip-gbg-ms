package serviceguiden

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/logging"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type ServiceGuidenClient interface {
	Badplatser(ctx context.Context) ([]Content, error)
}

type client struct {
	serviceUrl string
	badplatser []Content
	contents   []Content
}

func New(url, filePath string) ServiceGuidenClient {
	c, err := loadContentsFromFile(filePath)
	if err != nil {
		c = []Content{}
	}

	return &client{
		serviceUrl: url,
		contents:   c,
	}
}

func loadContentsFromFile(filePath string) (content []Content, err error) {
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &content)
	if err != nil {
		return
	}

	return
}

var tracer = otel.Tracer("integration-cip-gbg-ms/serviceguiden")

func (sgc client) Get(ctx context.Context) ([]Content, error) {
	var err error

	ctx, span := tracer.Start(ctx, "integration-cip-gbg-ms/serviceguiden/get")
	defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()

	httpClient := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sgc.serviceUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data from serviceguiden: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve data from serviceguiden, expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	contents := struct {
		Content []Content `json:"content"`
	}{}

	err = json.Unmarshal(body, &contents)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return contents.Content, err
}

func (sgc *client) Badplatser(ctx context.Context) ([]Content, error) {
	logger := logging.GetFromContext(ctx)

	if len(sgc.badplatser) > 0 {
		logger.Debug().Msgf("returning %d previously fetched beaches", len(sgc.badplatser))
		return sgc.badplatser, nil
	}

	if len(sgc.contents) == 0 {
		content, err := sgc.Get(ctx)
		if err != nil {
			return nil, err
		}
		sgc.contents = content
	}

	logger.Debug().Msgf("%d contents fetched from ServiceGuiden", len(sgc.contents))

	for _, c := range sgc.contents {
		if !c.Deleted {
			for _, st := range c.ServiceTypes {
				if st.Name == "Badplatser" {
					sgc.badplatser = append(sgc.badplatser, c)
				}
			}
		}
	}

	logger.Debug().Msgf("%d beaches found", len(sgc.badplatser))

	return sgc.badplatser, nil
}
