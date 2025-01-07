package dns

import "net/http"

type Api interface {
	UpdateRecord(ip string) error
}

type CloudFlareAPI struct {
	Token   string
	BaseURL string
	Client  *http.Client
}
