package ddns

import (
	"fmt"
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/miekg/dns"
)

type Server struct {
	// AllowedAPIKeys is a list of strings that will be allowed by the server
	// when provided as a bearer token in the authorization header.
	AllowedAPIKeys []string

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

func (s *Server) Set(domain string, ip net.IP) {
	if s.Domains == nil {
		s.Domains = map[string]net.IP{}
	}
	s.Domains[domain] = ip
}

func (s *Server) Listen() error {
	go s.listenHTTP(s.getHTTPListener())
	s.listenDNS(s.getDNSListener())
	return nil
}

func (s *Server) listenHTTP(listener string) error {
	http.HandleFunc("POST /api/v1/update", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%+v\n", s.Domains)
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !s.validateToken(token) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(nil)
			return
		}

		q := r.URL.Query()
		domain := q.Get("domain")
		ip := net.ParseIP(q.Get("ip"))

		if domain == "" || ip == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}

		// Only allow changing domains that already exist in the map
		if _, ok := s.Domains[domain]; !ok {
			w.WriteHeader(http.StatusForbidden)
			w.Write(nil)
			return
		}

		s.Domains[domain] = ip
		w.Write(nil)
	})

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
	return key != "" && slices.Contains(s.AllowedAPIKeys, key)
}

func (s *Server) getDNSListener() string {
	if s.DNSListener == "" {
		return ":5333"
	}
	return s.DNSListener
}
func (s *Server) getHTTPListener() string {
	if s.HTTPListener == "" {
		return ":8989"
	}
	return s.HTTPListener
}
