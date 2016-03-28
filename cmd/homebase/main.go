package main

import (
	"flag"
	"fmt"
	"github.com/bsdavidson/homebase"
	"log"
	"net/http"
)

var domainName, recordName, token string

func init() {
	flag.StringVar(&domainName, "domain", "", "Domain name hosted with Digital Ocean")
	flag.StringVar(&recordName, "record", "", "Subdomain to update")
	flag.StringVar(&token, "token", "", "Digital Ocean API token")
}

func main() {
	flag.Parse()
	if domainName == "" {
		log.Fatal("-domain argument is required")
	}
	if recordName == "" {
		log.Fatal("-record argument is required")
	}
	if token == "" {
		log.Fatal("-token argument is required")
	}
	client := &http.Client{}

	ip, err := homebase.GetPublicIP(client)
	if err != nil {
		log.Fatal("Error getting IP:", err)
	}

	record, err := homebase.GetRecordByName(client, domainName, recordName, token)
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
	fmt.Println(ip)
}
