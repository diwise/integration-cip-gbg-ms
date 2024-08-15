package serviceguiden

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestUnmarshalContent(t *testing.T) {
	is := is.New(t)
	var content Content
	err := json.Unmarshal([]byte(askimsbadet_json), &content)
	is.NoErr(err)
}

func TestInriktning(t *testing.T) {
	is := is.New(t)
	var content Content
	err := json.Unmarshal([]byte(askimsbadet_json), &content)
	is.NoErr(err)
	is.Equal("Hav", content.Inriktning())
}

func TestBeachTypes(t *testing.T) {
	is := is.New(t)
	var content Content
	err := json.Unmarshal([]byte(askimsbadet_json), &content)
	is.NoErr(err)
	is.Equal(1, len(content.BeachTypes()))
	is.Equal("Hav", content.BeachTypes()[0])
}

func TestAreaServed(t *testing.T) {
	is := is.New(t)
	var content Content
	err := json.Unmarshal([]byte(askimsbadet_json), &content)
	is.NoErr(err)
	is.Equal("Sydväst", content.AreaServed())
}

func TestIsBadplats(t *testing.T) {
	is := is.New(t)
	var content Content
	err := json.Unmarshal([]byte(askimsbadet_json), &content)
	is.NoErr(err)
	is.True(content.IsBadplats())
}

func TestUnmarshalServiceGuiden(t *testing.T) {
	is := is.New(t)
	f, err := os.Open("../../../../assets/test/serviceguiden_trim.json")
	is.NoErr(err)
	b, err := io.ReadAll(f)
	is.NoErr(err)
	var contents ServiceGuiden
	err = json.Unmarshal(b, &contents)
	is.NoErr(err)
	/*
		b,err = json.Marshal(contents)
		is.NoErr(err)
		err = os.WriteFile("../../../../assets/test/serviceguiden_trim.json", b, os.ModePerm)
		is.NoErr(err)
	*/
}

const askimsbadet_json string = `
{
    "id": "61e0a244cfc4d247cca95f4e",
    "name": "Askimsbadet",
    "siteUrl": "",
    "functions": [],
    "organization": {
        "id": "61e0a233cfc4d247cca954a9",
        "name": "Stadsmiljöförvaltningen",
        "externalId": "",
        "internalId": "N400"
    },
    "serviceTypes": [
        {
            "id": "61e0a232cfc4d247cca95394",
            "name": "Badplatser",
            "attributes": [
                {
                    "id": "61e18e4bc0e46817d17bb0fa",
                    "name": "Inriktning",
                    "values": [
                        {
                            "id": "61e18e4bc0e46817d17bb0f5",
                            "name": "Hav"
                        }
                    ]
                },
                {
                    "id": "64802d95658d671285c3347d",
                    "name": "Toalett",
                    "values": [
                        {
                            "id": "64802d95658d671285c3347a",
                            "name": "Toalett öppen under badsäsong"
                        }
                    ]
                },
                {
                    "id": "64802e4b658d671285c34a32",
                    "name": "Badservice",
                    "values": [
                        {
                            "id": "64802e4b658d671285c34a2d",
                            "name": "Sandstrand"
                        },
                        {
                            "id": "64802e4b658d671285c34a31",
                            "name": "Rullstolsramp ner i vattnet"
                        }
                    ]
                },
                {
                    "id": "64802c95658d671285c3171e",
                    "name": "Nakenbad",
                    "values": []
                },
                {
                    "id": "64802ea2658d671285c358ac",
                    "name": "Service",
                    "values": [
                        {
                            "id": "64802ea2658d671285c358ab",
                            "name": "Grillplats"
                        },
                        {
                            "id": "64802ea2658d671285c358aa",
                            "name": "Kiosk eller café"
                        }
                    ]
                },
                {
                    "id": "64802c53658d671285c308c6",
                    "name": "Hund tillåtet",
                    "values": [
                        {
                            "id": "6481e401e5b7910389a4efb3",
                            "name": "Nej"
                        }
                    ]
                }
            ]
        }
    ],
    "contacts": [
        {
            "id": "624f88ce3b546c0b5cd7113d",
            "name": "",
            "mobilePhone": {
                "e164": "",
                "display": ""
            },
            "phone": {
                "e164": "+46313650000",
                "display": "031-365 00 00"
            },
            "email": "",
            "contactCenter": true
        }
    ],
    "fax": {
        "e164": "",
        "display": ""
    },
    "businessId": 3683,
    "visitingAddress": "",
    "postalAddress": {
        "street": "",
        "city": "",
        "postalCode": ""
    },
    "position": {
        "latitude": 57.62595719307582,
        "longitude": 11.92624964921406
    },
    "description": "<p><strong>Avrådan från bad vid Askimsbadet</strong></p>\r\n<p><span>På grund av höga bakteriehalter har Göteborgs Stad beslutat om avrådan från bad vid Askimsbadet från och med 13 augusti. Avrådan från bad gäller till dess att vattenproverna visar att vattnet är tjänligt att bada i.</span></p>\r\n<p>En av Göteborgs mest populära badplatser med långgrund sandstrand och stora gräsytor för lek, spel och solbad. Här kan du sola och bada från den 259 meter långa badpiren, som är Sveriges längsta. Ramper ner till sandstranden och från piren gör det lättare för rullstolsburna att bada. Vid högsäsong finns badplatsvärdar på plats på Askimsbadet.</p>\r\n<p><strong>Service vid badet</strong></p>\r\n<p><a href=\"/wps/portal?uri=gbglnk%3a20201219207511\" target=\"_self\">Toaletterna är öppna under badsäsongen.</a></p>\r\n<p>På Askimsbadet är tre toaletter öppna fram till 10 juni, när samtliga toaletter öppnar.</p>\r\n<p>Utrustning och service: Kafé, lekplats, bangolf, utegym, utomhus- och inomhusduschar, toaletter, tillgänglighetsanpassad toalett och toalett med skötbord, omklädningsrum, beachvolleybollplan, surfcenter och grillplatser.</p>\r\n<p>På Askimsbadet finns även utställningsrummet Långgrundet som är öppet under badsäsongen.</p>\r\n<p><strong>Hundförbud på badet</strong></p>\r\n<p><span>Det är hundförbud på Askimsbadet mellan 1 maj och 15 september</span><span>.</span></p>\r\n<p>Från och med 30 juni får du passera Askimsbadet med kopplad hund. Det innebär att du nu får passera badplatsen via den grusade gången (Askims strandväg), ovanför lekplatsen, samt gångvägen till och från parkeringen.</p>\r\n<p>Tänk på att det fortfarande är förbjudet att passera badplatsen via gångvägen närmast stranden.</p>\r\n<p><strong>Närmaste hållplats:</strong> Askimsbadet. Från hållplatsen är det 290 meter till badplatsen.</p>\r\n<p><strong>Närmaste flexlinjemötesplats: </strong>1380 Askimsbadet ( Flexlinjen Askim och Flexlinjen Frölunda-Sisjön)</p>\r\n<p><a href=\"https://www.parkeringgoteborg.se/hitta-parkering/?searchtext=askimsbadet&amp;VisitOrRent=Visit&amp;parkingtype=1&amp;vehicletype=1&amp;SubmitSearchParking_ParkingPage_Visit=Visa\" target=\"_self\">Hitta närmaste parkering hos Parkering Göteborg</a></p>\r\n<p>Om papperskorgen eller containern är full, ta med dig skräpet hem. Tack för din hjälp!</p>",
    "districtOrganization": "Askim-Frölunda-Högsbo",
    "primaryArea": "",
    "cityArea": "Sydväst",
    "subCityArea": "Askim",
    "accessibilityUrl": "https://www.t-d.se/sv/TD2/Avtal/Goteborgs-stad/Askimsbadet/",
    "deleted": false,
    "images": [
        {
            "small": "https://s3.eu-north-1.amazonaws.com/gbg.serviceguiden/61e0a244cfc4d247cca95f4e_62138ca84c6152258898131f_small.jpg",
            "small2x": "https://s3.eu-north-1.amazonaws.com/gbg.serviceguiden/61e0a244cfc4d247cca95f4e_62138ca84c6152258898131f_small_2x.jpg",
            "medium": "https://s3.eu-north-1.amazonaws.com/gbg.serviceguiden/61e0a244cfc4d247cca95f4e_62138ca84c6152258898131f_medium.jpg",
            "medium2x": "https://s3.eu-north-1.amazonaws.com/gbg.serviceguiden/61e0a244cfc4d247cca95f4e_62138ca84c6152258898131f_medium_2x.jpg",
            "large": "https://s3.eu-north-1.amazonaws.com/gbg.serviceguiden/61e0a244cfc4d247cca95f4e_62138ca84c6152258898131f_large.jpg",
            "large2x": "https://s3.eu-north-1.amazonaws.com/gbg.serviceguiden/61e0a244cfc4d247cca95f4e_62138ca84c6152258898131f_large_2x.jpg",
            "altText": "Brygga med bland annat ramp ner i vattnet. ",
            "description": "Foto: Peter Svenson"
        }
    ],
    "eventLogs": [
        {
            "userName": "",
            "userId": "1",
            "eventAction": "UPDATED",
            "eventDateTime": "2022-04-07 10:43:14"
        }
    ]
}
`
