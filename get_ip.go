package homebase

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// HTTPGetter is an interface that implements http.Client's Get method so that
// mock objects can be passed in for tests.
type HTTPGetter interface {
	Get(url string) (resp *http.Response, err error)
}

// GetPublicIP returns the public IP address for the current machine.
func GetPublicIP(c HTTPGetter) (net.IP, error) {
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
