package middleware

import (
	"net/http"
	"regexp"
	"strings"
)

// Inspired by https://github.com/gorilla/handlers/blob/master/proxy_headers.go

var (
	// De-facto standard header keys.
	xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	//xForwardedHost   = http.CanonicalHeaderKey("X-Forwarded-Host")
	xForwardedProto  = http.CanonicalHeaderKey("X-Forwarded-Proto")
	xForwardedScheme = http.CanonicalHeaderKey("X-Forwarded-Scheme")
	xRealIP          = http.CanonicalHeaderKey("X-Real-IP")

	xEnvbinProxyChain = http.CanonicalHeaderKey("X-Envbin-Proxy-Chain")
)

var (
	// RFC7239 defines a new "Forwarded: " header designed to replace the
	// existing use of X-Forwarded-* headers.
	// e.g. Forwarded: for=192.0.2.60;proto=https;by=203.0.113.43
	forwarded      = http.CanonicalHeaderKey("Forwarded")
	forwardedRegex = regexp.MustCompile(`(?i)(?:for=)([^(;|,| )]+)(?:;proto=)?(https|http)?`)
)

func proxyHeaders(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fwd := getProxies(r)
		r.Header.Set(xEnvbinProxyChain, strings.Join(fwd, ", "))

		// Call the next handler in the chain.
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// getProxies retrieves the IP from the X-Forwarded-For, X-Real-IP and RFC7239
// Forwarded headers (in that order).
func getProxies(r *http.Request) []string {
	// TODO
	// * unit test
	// * integration test

	var addrs []string

	if fwd := r.Header.Get(xForwardedFor); fwd != "" {
		addrs = strings.Split(fwd, ", ")

		// only one of these, and it's for the first hop (client -> first proxy)
		var scheme string
		if proto := r.Header.Get(xForwardedProto); proto != "" {
			scheme = strings.ToLower(proto)
		} else if proto = r.Header.Get(xForwardedScheme); proto != "" {
			scheme = strings.ToLower(proto)
		}
		if scheme != "" {
			addrs[0] = addrs[0] + " (" + scheme + ")"
		}
	} else if fwd := r.Header.Get(xRealIP); fwd != "" {
		// X-Real-IP should only contain one IP address (the client making the
		// request).
		addrs = []string{fwd}
	} else if fwd := r.Header.Get(forwarded); fwd != "" {
		// format: "for=1.1.1.1;proto=https, for=4.4.4.4;proto=http"
		if matches := forwardedRegex.FindAllStringSubmatch(fwd, -1); len(matches) > 0 {
			// IPv6 addresses in Forwarded headers are quoted-strings. We strip
			// these quotes.
			addrs = []string{}
			for _, match := range matches {
				addr := strings.Trim(match[1], `"`)
				if match[2] != "" {
					addr = addr + " (" + match[2] + ")"
				}
				addrs = append(addrs, addr)
			}
		}
	}

	return addrs
}
