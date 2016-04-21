package homebase

import (
	"log"
	"net"
	"net/http"
)

// GetAndUpdate looks up the public IP address and updates a Digital Ocean
// subdomain with that IP.
func GetAndUpdate(domainName, recordName, token string) (net.IP, error) {
	client := &http.Client{}

	ip, err := GetPublicIP(client)
	if err != nil {
		log.Fatal("Error getting IP:", err)
	}

	record, err := GetRecordByName(client, domainName, recordName, token)
	if err != nil {
		log.Fatal("Error getting record: ", err)
	}

	if record.Type != "A" {
		log.Fatal("Record type must be A, was: ", record.Type)
	}

	record.Data = ip.String()
	err = record.Save(client, domainName, token)
	if err != nil {
		log.Fatal("Error saving record: ", err)
	}
	return ip, err
}
