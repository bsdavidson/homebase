package main

import (
	"flag"
	"fmt"
	"github.com/bsdavidson/homebase"
	"log"
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

	record, err := homebase.GetAndUpdate(domainName, recordName, token)
	if err != nil {
		log.Fatal("Error Getting and Updating: ", err)
	}

	fmt.Println(record)

}
