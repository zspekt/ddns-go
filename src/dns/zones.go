package dns

import (
	"log"

	"github.com/zspekt/ddns-go/pkg/utils"
)

func (api *CloudFlareAPI) getDNSZones() (*DNSZones, error) {
	var DNSZones *DNSZones = new(DNSZones)

	request, err := newRequestWithToken("GET", api.BaseURL+"/zones/", api.Token, nil)
	if err != nil {
		return nil, err
	}

	response, err := api.Client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	err = utils.DecodeJson(response.Body, DNSZones)
	if err != nil {
		log.Fatal("getZonesIDs", err)
	}
	err = checkForErrors(DNSZones.CommonResponse)
	return DNSZones, err
}
