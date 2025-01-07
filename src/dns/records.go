package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kr/pretty"

	"github.com/zspekt/ddns-go/pkg/utils"
)

// makes the put request to cloudflare's api to update all A records to point
// to the ip that's been passed in. returns any error found throughout the execution
// of the function, or any error from cloudflare's response
func (api *CloudFlareAPI) UpdateRecord(ip string) error {
	DNSZones, err := api.getDNSZones()
	if err != nil {
		log.Fatal(err) // TODO: don't crash
	}
	zoneID := DNSZones.Result[0].ID

	DNSRecords, err := api.getDNSRecords(zoneID)
	if err != nil {
		log.Fatal(err) // TODO: don't crash
	}

	filteredDNSRecords := filterRecords(DNSRecords.Result, filterRecordsWithType, "A")

	fmt.Println("final filtered dns records ahead === === === === === === === ")
	fmt.Println(pretty.Println(filteredDNSRecords))

	for _, d := range filteredDNSRecords {
		r, err := api.putRecord(d, ip)
		if err != nil {
			log.Fatal(err)
		}

		respp := new(DNSPutResponse)
		err = utils.DecodeJson(io.Reader(r.Body), respp)
		if err != nil {
			log.Fatal("error decoding final response")
		}
		fmt.Println(pretty.Println(respp))
	}
	return nil
}

func (api *CloudFlareAPI) getDNSRecords(zoneID string) (*DNSRecords, error) {
	var DNSRecords *DNSRecords = new(DNSRecords)

	request, err := newRequestWithToken(
		"GET",
		api.BaseURL+"/zones/"+zoneID+"/dns_records/",
		api.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	response, err := api.Client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	err = utils.DecodeJson(response.Body, DNSRecords)
	if err != nil {
		log.Fatal("getDNSRecords", err)
	}
	err = checkForErrors(DNSRecords.CommonResponse)
	return DNSRecords, err
}

func (api *CloudFlareAPI) putRecord(rec DNSRecord, ip string) (*http.Response, error) {
	d := DNSPutRequest{
		Comment:  "go ddns lol",
		Name:     rec.Name,
		Proxied:  rec.Proxied,
		Settings: rec.Settings,
		Tags:     rec.Tags,
		TTL:      rec.TTL,
		Content:  ip,
		Type:     rec.Type,
	}
	jsonPayload, err := json.Marshal(d)
	r := bytes.NewReader(jsonPayload)

	url := api.BaseURL + "/zones/" + rec.ZoneID + "/dns_records/" + rec.ID
	req, err := newRequestWithToken(
		"PUT",
		url,
		api.Token,
		r,
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "appliaction/json; charset=UTF-8")

	return api.Client.Do(req)
}

func filterRecords(rs []DNSRecord, f func(DNSRecord, string) bool, t string) []DNSRecord {
	ret := make([]DNSRecord, 0)
	for _, r := range rs {
		if f(r, t) {
			ret = append(ret, r)
		}
	}
	return ret
}

func filterRecordsWithType(rec DNSRecord, t string) bool {
	if rec.Type != t {
		return false
	}
	return true
}
