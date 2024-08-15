package serviceguiden

import (
	"strings"
)

type ServiceGuiden struct {
	Contents []Content `json:"content"`
}

type Content struct {
	ID_              string        `json:"id"`
	Name_            string        `json:"name"`
	SiteURL          string        `json:"siteUrl"`
	//Organization     Organization  `json:"organization"`
	ServiceTypes     []ServiceType `json:"serviceTypes"`
	Contacts         []Contact     `json:"contacts"`
	//Fax              ContactMethod `json:"fax"`
	BusinessID_      int           `json:"businessId"`
	VisitingAddress  string        `json:"visitingAddress"`
	//PostalAddress    PostalAddress `json:"postalAddress"`
	Position_        Position      `json:"position"`
	Description_     string        `json:"description"`
	//DistrictOrg      string        `json:"districtOrganization"`
	PrimaryArea      string        `json:"primaryArea"`
	CityArea         string        `json:"cityArea"`
	SubCityArea      string        `json:"subCityArea"`
	AccessibilityURL string        `json:"accessibilityUrl"`
	Deleted          bool          `json:"deleted"`
	//Images           []Image       `json:"images"`
}

/*
id (från service guide),
namn (name från service guide),
inriktning (name från values under attributes, serviceTypes från service guide med kommatecken mellan dem),
hemsida (siteUrl från service guide),
besoksAdress (visitingAddress från service guide),
*/

type Beach interface {
	ID() string
	Name() string
	Description() string
	WebSite() string
	Address() string
	Inriktning() string
	BeachTypes() []string
	AreaServed() string
	AccessibilityUrl() string
	Position() Position
	BusinessId() int
}

func (r Content) Description() string {
	return r.Description_
}

func (r Content) Name() string {
	return r.Name_
}
func (r Content) BusinessId() int {
	return r.BusinessID_
}
func (r Content) AccessibilityUrl() string {
	return r.AccessibilityURL
}
func (r Content) Position() Position {
	return r.Position_
}

func (r Content) ID() string {
	return r.ID_
}

func (r Content) WebSite() string {
	return r.SiteURL
}

func (r Content) Address() string {
	return r.VisitingAddress
}

func (r Content) Inriktning() string {
	attrs := make([]string, 0)
	for _, serviceType := range r.ServiceTypes {
		for _, attr := range serviceType.Attributes {
			name := strings.TrimSpace(attr.Name)
			if strings.EqualFold(name, "inriktning") {
				if len(attr.Values) > 0 {
					for _, v := range attr.Values {
						attrs = append(attrs, strings.TrimSpace(v.Name))
					}
				}
			}
		}
	}
	return strings.Join(attrs, ", ")
}

func (r Content) BeachTypes() []string {
	bt := strings.Split(r.Inriktning(), ",")
	if len(bt) == 0 {
		return nil
	}
	beachTypes := make([]string, 0)
	for _, b := range bt {
		beachTypes = append(beachTypes, strings.TrimSpace(b))
	}
	return beachTypes
}

func (r Content) AreaServed() string {
	if r.PrimaryArea != "" {
		return r.PrimaryArea
	}
	if r.CityArea != "" {
		return r.CityArea
	}
	if r.SubCityArea != "" {
		return r.SubCityArea
	}
	return ""
}

func (r Content) IsBadplats() bool {
	if r.Deleted {
		return false
	}
	for _, serviceType := range r.ServiceTypes {
		if strings.EqualFold(serviceType.Name, "Badplatser") {
			return true
		}
	}
	return false
}

type Organization struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ExternalID string `json:"externalId"`
	InternalID string `json:"internalId"`
}

type ServiceType struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Values []Value `json:"values"`
}

type Value struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Contact struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	MobilePhone   ContactMethod `json:"mobilePhone"`
	Phone         ContactMethod `json:"phone"`
	Email         string        `json:"email"`
	ContactCenter bool          `json:"contactCenter"`
}

type ContactMethod struct {
	E164    string `json:"e164"`
	Display string `json:"display"`
}

type PostalAddress struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	PostalCode string `json:"postalCode"`
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Image struct {
	Small       string `json:"small"`
	Small2x     string `json:"small2x"`
	Medium      string `json:"medium"`
	Medium2x    string `json:"medium2x"`
	Large       string `json:"large"`
	Large2x     string `json:"large2x"`
	AltText     string `json:"altText"`
	Description string `json:"description"`
}
