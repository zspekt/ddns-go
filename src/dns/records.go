package dns

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/zspekt/ddns-go/pkg/utils"
)

// makes the put request to cloudflare's api to update all A records to point
// to the ip that's been passed in. returns any error found throughout the execution
// of the function, or any error from cloudflare's response
func (api *CloudFlareAPI) UpdateRecord(ip string) error {
	Zones, err := api.getZones()
	if err != nil {
		log.Fatal(err) // TODO: don't crash
	}
	zoneID := Zones.Result[0].ID

	Records, err := api.getRecords(zoneID)
	if err != nil {
		log.Fatal(err) // TODO: don't crash
	}

	filteredRecords := filterRecords(Records.Result, filterRecordsWithType, "A")

	for _, d := range filteredRecords {
		r, err := api.putRecord(d, ip)
		if err != nil {
			log.Fatal(err)
		}

		respp := new(PutResponse)
		err = utils.DecodeJson(io.Reader(r.Body), respp)
		if err != nil {
			log.Fatal("error decoding final response")
		}
	}
	return nil
}

func (api *CloudFlareAPI) getRecords(zoneID string) (*Records, error) {
	var Records *Records = new(Records)

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

	err = utils.DecodeJson(response.Body, Records)
	if err != nil {
		log.Fatal("getDNSRecords", err)
	}
	err = checkForErrors(Records.CommonResponse)
	return Records, err
}

func (api *CloudFlareAPI) putRecord(rec Record, ip string) (*http.Response, error) {
	d := PutRequest{
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

func filterRecords(rs []Record, f func(Record, string) bool, t string) []Record {
	ret := make([]Record, 0)
	for _, r := range rs {
		if f(r, t) {
			ret = append(ret, r)
		}
	}
	return ret
}

func filterRecordsWithType(rec Record, t string) bool {
	if rec.Type != t {
		return false
	}
	return true
}
