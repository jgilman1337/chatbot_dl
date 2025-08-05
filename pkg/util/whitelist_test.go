package util

import (
	"net/url"
	"testing"
)

func TestWhitelist(t *testing.T) {
	pplxWhitelist := WhitelistEntry{
		Host:  "www.perplexity.ai",
		Paths: []string{"/search", "/rest/thread", "/cdn-cgi/challenge-platform"},
	}
	pplxCdnWhitelist := WhitelistEntry{
		Host:  "pplx-next-static-public.perplexity.ai",
		Paths: []string{"/_next/static/chunks", "/_next/static/css"},
		Exts:  []string{`\.js$`, `\.css$`},
	}

	tests := make(map[string]bool, 0)

	//Should accept
	tests["https://pplx-next-static-public.perplexity.ai/_next/static/css/b6786aa76fc37924.css"] = true
	tests["https://pplx-next-static-public.perplexity.ai/_next/static/css/4fc1e14b7c7dfade.css"] = true
	tests["https://pplx-next-static-public.perplexity.ai/_next/static/css/37702e1215bdd867.css"] = true
	tests["https://www.perplexity.ai/search/ai-amplifies-false-memories-9iZN5JuFT5.9asR1Ntf._A"] = true
	tests["https://www.perplexity.ai/rest/thread/ai-amplifies-false-memories-9iZN5JuFT5.9asR1Ntf._A"] = true
	tests["https://pplx-next-static-public.perplexity.ai/_next/static/chunks/46324-1d5ecf9952fa20ea.js"] = true
	tests["https://pplx-next-static-public.perplexity.ai/_next/static/chunks/36020-157fbc80db2af79e.js"] = true
	tests["https://pplx-next-static-public.perplexity.ai/_next/static/chunks/62835-833f199d9949b8c8.js"] = true

	//Should reject
	tests["https://www.perplexity.ai/rest/event/analytics"] = false
	tests["https://browser-intake-datadoghq.com/api/v2/rum"] = false
	tests["https://accounts.google.com/gsi/client"] = false
	tests["https://static.cloudflareinsights.com/beacon.min.js/vcd15cbe7772f49c399c6a5babf22c1241717689176015"] = false

	//Run tests
	testRunner([]WhitelistEntry{pplxCdnWhitelist, pplxWhitelist}, tests, t)
}

func testRunner(whitelists []WhitelistEntry, tests map[string]bool, t *testing.T) {
	i := 1
	for turl, expected := range tests {
		//Parse the URL first
		u, err := url.Parse(turl)
		if err != nil {
			t.Fatalf("[url parse error] %s\n", err)
		}

		//Get the hostname and path
		hostname := u.Hostname()
		path := u.EscapedPath()

		actual := IsWhitelisted(hostname, path, whitelists...)
		if actual != expected {
			t.Fatalf("[test %d/%d failed] URL %s; got: %v, expected: %v", i, len(tests), u, actual, expected)
		}
		i++
	}
}
