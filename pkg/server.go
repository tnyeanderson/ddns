package ddns

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

const (
	DefaultDNSListener  = ":53"
	DefaultHTTPListener = ":3345"
)

type Server struct {
	// AllowedAPIKeys is a list of strings that will be allowed by the server
	// when provided as a bearer token in the authorization header. A nil matcher
	// allows changing the DNS record for any domain (no restrictions).
	AllowedAPIKeys map[string]*regexp.Regexp

	// HTTPListener is the TCP address that will be passed to
	// [http.ListenAndServe()].
	HTTPListener string

	// DNSListener is the network address that the DNS server will listen on. See
	// [dns.Server].
	DNSListener string

	// Domains is the map of a domain to an IP that will be returned by the DNS
	// server.
	Domains map[string]net.IP
}

// Allow is a convenience function for adding API keys which are allowed to
// change the DNS record for the domains matched by domainMatcher. Further
// needs should be handled via [Server.AllowedAPIKeys] directly.  A nil
// domainMatcher allows changing the DNS record for any domain (no
// restrictions).
func (s *Server) Allow(apiKey string, domainMatcher *regexp.Regexp) {
	if s.AllowedAPIKeys == nil {
		s.AllowedAPIKeys = map[string]*regexp.Regexp{}
	}
	s.AllowedAPIKeys[apiKey] = domainMatcher
}

// Set updates the DNS record for the provided domain.
func (s *Server) Set(domain string, ip net.IP) {
	if s.Domains == nil {
		s.Domains = map[string]net.IP{}
	}
	s.Domains[domain] = ip
}

// Listen starts a DNS server and an HTTP server for the API, and blocks until
// either of them exits.
func (s *Server) Listen() error {
	wg := sync.WaitGroup{}
	// If one exits, end the program
	wg.Add(1)

	var out error

	go func() {
		if err := s.listenHTTP(s.getHTTPListener()); err != nil {
			slog.Error(err.Error())
			out = err
		}
		wg.Done()
	}()

	go func() {
		if err := s.listenDNS(s.getDNSListener()); err != nil {
			slog.Error(err.Error())
			out = err
		}
		wg.Done()
	}()

	wg.Wait()
	return out
}

func (s *Server) handleGetIP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, err := getCallerIP(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(ip.String()))
	}
}

func (s *Server) handleUpdateIP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate token
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !s.validateToken(token) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Validate domain
		domain := r.URL.Query().Get("domain")
		if domain == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Only allow changing domains that are allowed by the token
		if r := s.AllowedAPIKeys[token]; r != nil && !r.MatchString(domain) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Determine IP
		ipStr := r.URL.Query().Get("ip")
		if ipStr == "auto" {
			callerIP, err := getCallerIP(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			ipStr = callerIP.String()
		}

		// Validate IP
		ip := net.ParseIP(ipStr)
		if ip == nil || ip.IsLoopback() || ip.IsMulticast() || ip.IsUnspecified() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.Set(domain, ip)
		w.Write(nil)
	}
}

func (s *Server) listenHTTP(listener string) error {
	http.HandleFunc("GET /api/v1/ip", s.handleGetIP())
	http.HandleFunc("POST /api/v1/update", s.handleUpdateIP())
	return http.ListenAndServe(listener, nil)
}

func (s *Server) listenDNS(listener string) error {
	dnsServer := &dns.Server{Addr: listener, Net: "udp"}
	dns.HandleFunc(".", func(w dns.ResponseWriter, m *dns.Msg) {
		r := new(dns.Msg)
		r.SetReply(m)
		defer w.WriteMsg(r)
		for _, q := range r.Question {
			domain := strings.TrimSuffix(q.Name, ".")
			ip := s.Domains[domain]
			if q.Qtype == dns.TypeA && ip != nil {
				answer := &dns.A{
					A: ip,
					Hdr: dns.RR_Header{
						Name:   q.Name,
						Rrtype: q.Qtype,
						Class:  q.Qclass,
					},
				}
				r.Answer = append(r.Answer, answer)
			}
		}
	})
	return dnsServer.ListenAndServe()
}

func (s *Server) validateToken(key string) bool {
	if key == "" {
		return false
	}

	for k, _ := range s.AllowedAPIKeys {
		if k == key {
			return true
		}
	}

	return false
}

func (s *Server) getDNSListener() string {
	if s.DNSListener == "" {
		return DefaultDNSListener
	}
	return s.DNSListener
}
func (s *Server) getHTTPListener() string {
	if s.HTTPListener == "" {
		return DefaultHTTPListener
	}
	return s.HTTPListener
}

func getCallerIP(r *http.Request) (net.IP, error) {
	// Use X-Real-Ip header if available
	ip := r.Header.Get("X-Real-Ip")

	// Use X-Forwarded-For header if available
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}

	// Use request.RemoteAddr
	if ip == "" {
		lastColon := strings.LastIndex(r.RemoteAddr, ":")
		ip = r.RemoteAddr[:lastColon]
		ip = strings.Trim(ip, "[]")
	}

	// Ensure it is a valid IP
	out := net.ParseIP(ip)
	if out == nil {
		return nil, fmt.Errorf("not a valid ip: %s", ip)
	}

	return out, nil
}
