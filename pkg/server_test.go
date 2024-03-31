package ddns

import (
	"net"
	"regexp"
	"testing"
)

func TestServerAllow(t *testing.T) {
	s := Server{}
	perms := map[string]*regexp.Regexp{
		"mykey":     regexp.MustCompile(".*"),
		"mynullkey": nil,
		"thirdkey":  regexp.MustCompile(`^domain\.com$`),
	}
	for key, matcher := range perms {
		s.Allow(key, matcher)
	}
	if s.AllowedAPIKeys == nil {
		t.Fatalf("not initialized: Server.AllowedAPIKeys")
	}
	for key, matcher := range perms {
		got, ok := s.AllowedAPIKeys[key]
		if !ok {
			t.Fatalf("key missing from Server.AllowedAPIKeys: %s", key)
		}
		if got != matcher {
			t.Fatalf(`incorrect value for Server.AllowedAPIKeys["%s"], got: "%s", expected: "%s"`, key, got, matcher)
		}
	}

	// overwrite a single value
	key, matcher := "mykey", regexp.MustCompile("haha")
	s.Allow(key, matcher)
	got, ok := s.AllowedAPIKeys[key]
	if !ok {
		t.Fatalf("key missing from Server.AllowedAPIKeys: %s", key)
	}
	if got != matcher {
		t.Fatalf(`incorrect overwritten value for Server.AllowedAPIKeys["%s"], got: "%s", expected: "%s"`, key, got, matcher)
	}
}

func TestServerSet(t *testing.T) {
	s := Server{}
	domains := map[string]net.IP{
		"mydomain.com":      net.ParseIP("1.2.3.4"),
		"myotherdomain.com": net.ParseIP("4.3.2.1"),
	}
	for k, v := range domains {
		s.Set(k, v)
	}
	if s.Domains == nil {
		t.Fatalf("not initialized: Server.Domains")
	}
	for domain, ip := range domains {
		got, ok := s.Domains[domain]
		if !ok {
			t.Fatalf("key missing from Server.Domains: %s", domain)
		}
		if !got.Equal(ip) {
			t.Fatalf(`incorrect value for Server.Domains["%s"], got: "%s", expected: "%s"`, domain, got, ip)
		}
	}

	// overwrite a single value
	domain, ip := "myotherdomain", net.ParseIP("3.3.3.3")
	s.Set(domain, ip)
	got, ok := s.Domains[domain]
	if !ok {
		t.Fatalf("key missing from Server.Domains: %s", domain)
	}
	if !got.Equal(ip) {
		t.Fatalf(`incorrect overwritten value for Server.Domains["%s"], got: "%s", expected: "%s"`, domain, got, ip)
	}
}

func TestServerLoad(t *testing.T) {
	s := Server{
		HostsFile: "testdata/hosts.yaml",
	}
	if err := s.Load(); err != nil {
		t.Fatal(err)
	}
	domains := map[string]net.IP{
		"mydomain.com":      net.ParseIP("1.2.3.4"),
		"myotherdomain.com": net.ParseIP("4.3.2.1"),
	}
	if s.Domains == nil {
		t.Fatalf("not initialized: Server.Domains")
	}
	for domain, ip := range domains {
		got, ok := s.Domains[domain]
		if !ok {
			t.Fatalf("key missing from Server.Domains: %s", domain)
		}
		if !got.Equal(ip) {
			t.Fatalf(`incorrect value for Server.Domains["%s"], got: "%s", expected: "%s"`, domain, got, ip)
		}
	}

	// overwrite a single value
	//s.HostsFile = "testdata/hosts-new.yaml"
	//if err := s.Load(); err != nil {
	//	t.Fatal(err)
	//}
	//domain, ip := "mydomain.com", net.ParseIP("3.3.3.3")
	//got, ok := s.Domains[domain]
	//if !ok {
	//	t.Fatalf("key missing from Server.Domains: %s", domain)
	//}
	//if !got.Equal(ip) {
	//	t.Fatalf(`incorrect overwritten value for Server.Domains["%s"], got: "%s", expected: "%s"`, domain, got, ip)
	//}

	s.HostsFile = "testdata/hosts-new.yaml"
	if err := s.Load(); err != nil {
		t.Fatal(err)
	}
	// what got overwritten
	domains["mydomain.com"] = net.ParseIP("3.3.3.3")
	for domain, ip := range domains {
		got, ok := s.Domains[domain]
		if !ok {
			t.Fatalf("key missing from Server.Domains: %s", domain)
		}
		if !got.Equal(ip) {
			t.Fatalf(`incorrect value for Server.Domains["%s"], got: "%s", expected: "%s"`, domain, got, ip)
		}
	}

}
