package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain - hasMX - hasSPF - SPF Record - hasDMARC - DMARC record\n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
		fmt.Printf("\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error: could not read from input: %v\n", err)
	}	
}

func checkDomain(domain string){
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string
	
	mxRecords, err := net.LookupMX(domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spfi"){
			hasSPF = true
			spfRecord = record
			break
		}
	}

	 dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	 for _, record := range dmarcRecords {
		 if strings.HasPrefix(record, "v=DMARC1") {
			 hasDMARC = true
			 dmarcRecord = record
			 break
		 }
	 }

	 fmt.Printf("%v - %v - %v - %v - %v - %v", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}