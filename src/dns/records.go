package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kr/pretty"

	"github.com/zspekt/ddns-go/pkg/utils"
)

// makes the put request to cloudflare's api to update all A records to point
// to the ip that's been passed in. returns any error found throughout the execution
// of the function, or any error from cloudflare's response
func UpdateDnsRecord(ip, token string) error {
	client := http.Client{Timeout: 10 * time.Second}
	var baseURL string = "https://api.cloudflare.com/client/v4"

	DNSZones, err := getDNSZones(&client, baseURL, token)
	if err != nil {
		log.Fatal(err) // TODO: don't crash
	}
	zoneID := DNSZones.Result[0].ID

	DNSRecords, err := getDNSRecords(&client, baseURL, zoneID, token)
	if err != nil {
		log.Fatal(err) // TODO: don't crash
	}

	filteredDNSRecords := filterDNSRecords(DNSRecords.Result, "A", filterDNSRecsWithType)

	fmt.Println("final filtered dns records ahead === === === === === === === ")
	fmt.Println(pretty.Println(filteredDNSRecords))

	for _, d := range filteredDNSRecords {
		r, err := dnsPutRequest(
			&client,
			d,
			ip,
			token,
		)
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

func getDNSRecords(client *http.Client, baseURL, zoneID, token string) (*DNSRecords, error) {
	var DNSRecords *DNSRecords = new(DNSRecords)

	request, err := newRequestWithToken("GET", baseURL+"/zones/"+zoneID+"/dns_records/", token, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
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

func dnsPutRequest(client *http.Client, rec DNSRecord, ip, token string) (*http.Response, error) {
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

	url := "https://api.cloudflare.com/client/v4/zones/" + rec.ZoneID + "/dns_records/" + rec.ID
	req, err := newRequestWithToken(
		"PUT",
		url,
		token,
		r,
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "appliaction/json; charset=UTF-8")

	return client.Do(req)
}

func filterDNSRecords(rs []DNSRecord, t string, f func(DNSRecord, string) bool) []DNSRecord {
	ret := make([]DNSRecord, 0)
	for _, r := range rs {
		if f(r, t) {
			ret = append(ret, r)
		}
	}
	return ret
}

func filterDNSRecsWithType(rec DNSRecord, t string) bool {
	if rec.Type != t {
		return false
	}
	return true
}
