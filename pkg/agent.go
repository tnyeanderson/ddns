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
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(nil))
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (a *Agent) UpdateIP(domain, ip string) error {
	server := strings.TrimSuffix(a.ServerAddress, "/")
	url := fmt.Sprintf("%s/api/v1/update?domain=%s&ip=%s", server, domain, ip)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(nil))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.APIKey))
	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
