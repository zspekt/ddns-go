package dns

import (
	"net/http"
	"time"
)

type CloudFlareAPI struct {
	Token   string
	BaseURL string
	Client  *http.Client
}

type Record struct {
	Comment    any       `json:"comment"`
	Content    string    `json:"content"`
	CreatedOn  time.Time `json:"created_on"`
	ID         string    `json:"id"`
	Meta       Meta      `json:"meta"`
	ModifiedOn time.Time `json:"modified_on"`
	Name       string    `json:"name"`
	Proxiable  bool      `json:"proxiable"`
	Proxied    bool      `json:"proxied"`
	Settings   struct{}  `json:"settings"`
	Tags       []any     `json:"tags"`
	TTL        int       `json:"ttl"`
	Type       string    `json:"type"`
	ZoneID     string    `json:"zone_id"`
	ZoneName   string    `json:"zone_name"`
}

type Zone struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Status              string    `json:"status"`
	Paused              bool      `json:"paused"`
	Type                string    `json:"type"`
	DevelopmentMode     int       `json:"development_mode"`
	NameServers         []string  `json:"name_servers"`
	OriginalNameServers []string  `json:"original_name_servers"`
	OriginalRegistrar   string    `json:"original_registrar"`
	OriginalDnshost     any       `json:"original_dnshost"`
	ModifiedOn          time.Time `json:"modified_on"`
	CreatedOn           time.Time `json:"created_on"`
	ActivatedOn         time.Time `json:"activated_on"`
	Meta                struct {
		Step                   int  `json:"step"`
		CustomCertificateQuota int  `json:"custom_certificate_quota"`
		PageRuleQuota          int  `json:"page_rule_quota"`
		PhishingDetected       bool `json:"phishing_detected"`
	} `json:"meta"`
	Owner struct {
		ID    any    `json:"id"`
		Type  string `json:"type"`
		Email any    `json:"email"`
	} `json:"owner"`
	Account struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Tenant struct {
		ID   any `json:"id"`
		Name any `json:"name"`
	} `json:"tenant"`
	TenantUnit struct {
		ID any `json:"id"`
	} `json:"tenant_unit"`
	Permissions []string `json:"permissions"`
	Plan        struct {
		ID                string `json:"id"`
		Name              string `json:"name"`
		Price             int    `json:"price"`
		Currency          string `json:"currency"`
		Frequency         string `json:"frequency"`
		IsSubscribed      bool   `json:"is_subscribed"`
		CanSubscribe      bool   `json:"can_subscribe"`
		LegacyID          string `json:"legacy_id"`
		LegacyDiscount    bool   `json:"legacy_discount"`
		ExternallyManaged bool   `json:"externally_managed"`
	} `json:"plan"`
}

type Zones struct {
	Result     []Zone     `json:"result"`
	ResultInfo ResultInfo `json:"result_info"`
	CommonResponse
}

type Records struct {
	Result     []Record   `json:"result"`
	ResultInfo ResultInfo `json:"result_info"`
	CommonResponse
}

type PutRequest struct {
	Comment  string   `json:"comment"`
	Name     string   `json:"name"`
	Proxied  bool     `json:"proxied"`
	Settings struct{} `json:"settings"`
	Tags     []any    `json:"tags"`
	TTL      int      `json:"ttl"`
	Content  string   `json:"content"`
	Type     string   `json:"type"`
}

type PutResponse struct {
	Result any `json:"result"`
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Messages []any `json:"messages"`
	Success  bool  `json:"success"`
}

type Meta struct {
	AutoAdded           bool `json:"auto_added"`
	ManagedByApps       bool `json:"managed_by_apps"`
	ManagedByArgoTunnel bool `json:"managed_by_argo_tunnel"`
}

type ResultInfo struct {
	Count      int `json:"count"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

type CommonResponse struct {
	Errors   []Error `json:"errors"`
	Messages []any   `json:"messages"`
	Success  bool    `json:"success"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
