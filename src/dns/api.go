package dns

import "net/http"

type Api interface {
	getZones() (*Zones, error)
	UpdateRecord(ip string) error
	getRecords(zoneID string) (*Records, error)
	putRecord(rec Record, ip string) (*http.Response, error)
}

type CloudFlareAPI struct {
	Token   string
	BaseURL string
	Client  *http.Client
}
