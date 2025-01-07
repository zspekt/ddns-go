package dns

import (
	"log"

	"github.com/zspekt/ddns-go/pkg/utils"
)

func (api *CloudFlareAPI) getZones() (*Zones, error) {
	var Zones *Zones = new(Zones)

	request, err := newRequestWithToken("GET", api.BaseURL+"/zones/", api.Token, nil)
	if err != nil {
		return nil, err
	}

	response, err := api.Client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	err = utils.DecodeJson(response.Body, Zones)
	if err != nil {
		log.Fatal("getZonesIDs", err)
	}
	err = checkForErrors(Zones.CommonResponse)
	return Zones, err
}
