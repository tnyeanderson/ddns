package ddns

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const DefaultServerAddress = "http://localhost:3345"

type Agent struct {
	// ServerAddress is the scheme/host/port of the DDNS API server, not
	// including the /api base path. If not set, [DefaultServerAddress] will be
	// used.
	ServerAddress string

	// APIKey will be used to authenticate to the DDNS API.
	APIKey string
}

// DetermineIP returns the public IP of the caller as seen by the DDNS API
// server. Uses the /api/v1/ip endpoint.
func (a *Agent) DetermineIP() (string, error) {
	url := fmt.Sprintf("%s/api/v1/ip", a.getServerAddress())
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET request to %s returned unexpected status code: %d", url, res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// UpdateIP updates the IP for a given DDNS domain. Uses the /api/v1/update
// endpoint.
func (a *Agent) UpdateIP(domain, ip string) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/update?domain=%s&ip=%s", a.getServerAddress(), domain, ip)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(nil))
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.APIKey))
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	updated := false
	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
		updated = true
	default:
		return false, fmt.Errorf("POST request to %s returned unexpected status code: %d", url, res.StatusCode)
	}

	return updated, nil
}

func (a *Agent) getServerAddress() string {
	server := DefaultServerAddress
	if a.ServerAddress != "" {
		server = a.ServerAddress
	}
	return strings.TrimSuffix(server, "/")
}
