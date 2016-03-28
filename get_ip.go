package homebase

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type HttpGetter interface {
	Get(url string) (resp *http.Response, err error)
}

func GetPublicIP(c HttpGetter) (net.IP, error) {
	resp, err := c.Get("http://checkip.amazonaws.com")
	if err != nil {
		return nil, fmt.Errorf("Error requesting IP: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Bad status requesting IP: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading body: %v", err)
	}
	ip := net.ParseIP(strings.TrimSpace(string(body)))
	if ip == nil {
		return nil, fmt.Errorf("Invalid IP: %s", string(body))
	}

	return ip, nil
}
