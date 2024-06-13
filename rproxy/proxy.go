package rproxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

type DynamicProxyServer struct {
	RuleProvider DynamicProxyProvider
}

func NewDynamicProxyServer() *DynamicProxyServer {
	dps := &DynamicProxyServer{
		RuleProvider: NewInmemoryProxyProvider(),
	}
	dps.RuleProvider.AddRule(&MatchRule{
		MatchHost: "fakeserver.com",
		MatchPath: "/v1",
	}, &Target{
		Host:        "localhost:8000",
		RewritePath: "/",
	})
	dps.RuleProvider.AddRule(&MatchRule{
		MatchHost: "fakeserver.com",
		MatchPath: "/v2",
	}, &Target{
		Host:        "localhost:8100",
		RewritePath: "/",
	})
	return dps
}

func (d *DynamicProxyServer) NeedDispatch(c *gin.Context) bool {
	rule := d.RuleProvider.MatchRule(c)
	if rule == nil {
		return false
	}
	if rule.Target == nil {
		return false
	}
	p := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = rule.Target.Host
			req.URL.Path = rule.Target.RewritePath
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
			req.Host = rule.Target.Host
		},
	}
	p.ServeHTTP(c.Writer, c.Request)
	return true
}

func (d *DynamicProxyServer) Dispatch(c *gin.Context) {
	if d.NeedDispatch(c) {
		c.Abort()
		return
	}
	c.Next()
}

type DynamicProxyProvider interface {
	AddRule(rule *MatchRule, target *Target)
	MatchRule(c *gin.Context) *Rule
}

type MatchRule struct {
	MatchHost   string
	MatchPath   string
	MatchHeader map[string]string
	MatchQuery  map[string]string
}

type Target struct {
	Host        string
	RewritePath string
}

type Rule struct {
	MatchRule *MatchRule
	Target    *Target
	// TODO: add target auth related rule
}

type InmemoryProxyProvider struct {
	rules []*Rule
}

func NewInmemoryProxyProvider() *InmemoryProxyProvider {
	return &InmemoryProxyProvider{
		rules: make([]*Rule, 0),
	}
}

func (i *InmemoryProxyProvider) AddRule(rule *MatchRule, target *Target) {
	i.rules = append(i.rules, &Rule{MatchRule: rule, Target: target})
}

func (i *InmemoryProxyProvider) MatchRule(c *gin.Context) *Rule {
	for _, rule := range i.rules {
		println(rule.MatchRule.MatchHost, c.Request.Host)
		if rule.MatchRule.MatchHost == c.Request.Host && rule.MatchRule.MatchPath == c.Request.URL.Path {
			return rule
		}
	}
	return nil
}
