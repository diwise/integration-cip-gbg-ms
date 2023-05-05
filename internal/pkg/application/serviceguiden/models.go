package serviceguiden

import (
	"strings"
)

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

func (c Content) BeachTypes() []string {
	facilities := strings.Split(c.Inriktning(), ",")

	if len(facilities) == 0 {
		return nil
	}

	s := []string{""}
	for _, f := range facilities {
		s = append(s, strings.TrimSpace(f))
	}

	return s
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
