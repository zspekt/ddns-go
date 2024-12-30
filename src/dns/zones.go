package dns

import (
	"log"
	"net/http"

	"github.com/zspekt/ddns-go/pkg/utils"
)

func getDNSZones(client *http.Client, baseURL, token string) (*DNSZones, error) {
	var DNSZones *DNSZones = new(DNSZones)

	request, err := newRequestWithToken("GET", baseURL+"/zones/", token, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
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
