package ddns

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Agent struct {
	ServerAddress string
	APIKey        string
}

func (a *Agent) DetermineIP() (string, error) {
	server := strings.TrimSuffix(a.ServerAddress, "/")
	url := fmt.Sprintf("%s/api/v1/ip", server)
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

func (a *Agent) UpdateIP(domain, ip string) (bool, error) {
	server := strings.TrimSuffix(a.ServerAddress, "/")
	url := fmt.Sprintf("%s/api/v1/update?domain=%s&ip=%s", server, domain, ip)
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
