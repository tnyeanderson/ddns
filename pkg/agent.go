package ddns

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type Agent struct {
	ServerAddress string
	APIKey        string
}

func (a *Agent) MyIP() string {
	// TODO: Get public IP
	return "1.2.3.4"
}

func (a *Agent) UpdateIP(domain, ip string) error {
	server := strings.TrimSuffix(a.ServerAddress, "/")
	url := fmt.Sprintf("%s/api/v1/update?domain=%s&ip=%s", server, domain, ip)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(nil))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.APIKey))
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
