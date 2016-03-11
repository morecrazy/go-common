package protocol

type RR struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Weight   int    `json:"weight,omitempty"`
	Text     string `json:"text,omitempty"`
	Mail     bool   `json:"mail,omitempty"` // Be an MX record. Priority becomes Preference.
	Ttl      uint32 `json:"ttl,omitempty"`
}

type SetDnsReq struct {
	URL string `json:"url"`
	RRs []RR   `json:"rrs"`
}

type DelDnsReq struct {
	URL string `json:"url"`
}

type GetDnsReq struct {
	URL string `json:"url"`
}

type GetDnsResp struct {
	URL string `json:"url"`
	RRs []RR   `json:"rrs"`
}
