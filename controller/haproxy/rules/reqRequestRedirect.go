package rules

import (
	"fmt"

	"github.com/haproxytech/models/v2"

	"github.com/haproxytech/kubernetes-ingress/controller/haproxy"
	"github.com/haproxytech/kubernetes-ingress/controller/haproxy/api"
	"github.com/haproxytech/kubernetes-ingress/controller/utils"
)

type RequestRedirect struct {
	RedirectCode int64
	RedirectPort int
	Host         string
	SSLRequest   bool
	SSLRedirect  bool
}

func (r RequestRedirect) GetType() haproxy.RuleType {
	return haproxy.REQ_REQUEST_REDIRECT
}

func (r RequestRedirect) Create(client api.HAProxyClient, frontend *models.Frontend, ingressACL string) error {
	if frontend.Mode == "tcp" {
		return fmt.Errorf("request redirection cannot be configured in TCP mode")
	}
	var rule string
	if r.SSLRedirect {
		rule = fmt.Sprintf("https://%%[hdr(host),field(1,:)]:%d%%[capture.req.uri]", r.RedirectPort)
	} else {
		scheme := "http"
		if r.SSLRequest {
			scheme = "https"
		}
		rule = fmt.Sprintf(scheme+"://%s%%[capture.req.uri]", r.Host)
	}
	httpRule := models.HTTPRequestRule{
		Index:      utils.PtrInt64(0),
		Type:       "redirect",
		RedirCode:  utils.PtrInt64(r.RedirectCode),
		RedirValue: rule,
		RedirType:  "location",
	}
	return client.FrontendHTTPRequestRuleCreate(frontend.Name, httpRule, ingressACL)
}