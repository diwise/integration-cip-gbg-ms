package serviceguiden

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/logging"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type ServiceGuidenClient interface {
	Get(ctx context.Context) (*sgResponse, error)
	Badplatser(ctx context.Context) ([]Content, error)
}

type serviceGuidenClient struct {
	serviceUrl string
	badplatser []Content
}

type sgResponse struct {
	Content []Content `json:"content"`
}

type Content struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Organization struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"organization"`
	ServiceTypes []struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		Attributes []struct {
			Name   string `json:"name"`
			Values []struct {
				Name string `json:"name"`
			} `json:"values"`
		} `json:"attributes"`
	} `json:"serviceTypes"`
	Position struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"position"`
	PrimaryArea      string `json:"primaryArea"`
	CityArea         string `json:"cityArea"`
	SubCityArea      string `json:"subCityArea"`
	SiteUrl          string `json:"siteUrl"`
	AccessibilityUrl string `json:"accessibilityUrl"`
	Address          string `json:"visitingAddress"`
	PostalAddress    struct {
		Street     string `json:"street"`
		City       string `json:"city"`
		PostalCode string `json:"postalCode"`
	} `json:"postalAddress"`
	Deleted              bool   `json:"deleted"`
	Description          string `json:"description"`
	DistrictOrganization string `json:"districtOrganization"`
	BusinessId           int64  `json:"businessId"`
}

/*
id (från service guide),
namn (name från service guide),
inriktning (name från values under attributes, serviceTypes från service guide med kommatecken mellan dem),
hemsida (siteUrl från service guide),
besoksAdress (visitingAddress från service guide),
*/

func (c Content) Inriktning() string {
	str := ""
	for _, st := range c.ServiceTypes {
		for _, attr := range st.Attributes {
			n := strings.TrimSpace(attr.Name)
			if strings.EqualFold(n, "inriktning") {
				if len(attr.Values) > 0 {
					for _, v := range attr.Values {
						str = str + v.Name + ", "
					}
				}

				if str == "" {
					return ""
				}
				if strings.LastIndex(str, ",") == -1 {
					return str
				}

				return str[:strings.LastIndex(str, ",")]
			}
		}
	}

	return ""
}

func (c Content) Facilities() []string {
	var s []string
	facilities := strings.Split(c.Inriktning(), ",")

	if len(facilities) == 0 {
		return s
	}

	for _, f := range facilities {
		s = append(s, strings.TrimSpace(f))
	}

	return s
}

func (c Content) SeeAlso() string {
	return fmt.Sprintf("https://goteborg.se/wps/portal/start/kultur-och-fritid/fritid-och-natur/friluftsliv-natur-och/badplatser--utomhusbad/badplatser-utomhusbad/?id=%d", c.BusinessId)
}

func (c Content) AreaServed() string {
	if c.PrimaryArea != "" {
		return c.PrimaryArea
	} else if c.CityArea != "" {
		return c.CityArea
	} else if c.SubCityArea != "" {
		return c.SubCityArea
	}
	return ""
}

func New(url string) ServiceGuidenClient {

	return &serviceGuidenClient{
		serviceUrl: url,
	}
}

var tracer = otel.Tracer("integration-cip-gbg-ms/serviceguiden")

func (sgc serviceGuidenClient) Get(ctx context.Context) (*sgResponse, error) {
	var err error

	ctx, span := tracer.Start(ctx, "integration-cip-gbg-ms/serviceguiden/get")
	defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()

	log := logging.GetFromContext(ctx)

	httpClient := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sgc.serviceUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve data from serviceguiden")
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("failed to retrieve data from serviceguiden, expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		return nil, fmt.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return nil, err
	}

	var m sgResponse
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal model")
	}

	return &m, err
}

func (sgc *serviceGuidenClient) Badplatser(ctx context.Context) ([]Content, error) {
	if len(sgc.badplatser) > 0 {
		return sgc.badplatser, nil
	}

	resp, err := sgc.Get(ctx)
	if err != nil {
		return nil, err
	}

	for _, c := range resp.Content {
		if !c.Deleted {
			for _, st := range c.ServiceTypes {
				if st.Name == "Badplatser" {
					sgc.badplatser = append(sgc.badplatser, c)
				}
			}
		}
	}

	return sgc.badplatser, nil
}
