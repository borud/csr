// package main contains a utility to ensure you can parse CSR.
package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage %s <CSR file>", os.Args[0])
	}
	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error opening file %s: %v", filename, err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("error reading file %s: %v", filename, err)
	}

	pemBlock, _ := pem.Decode(data)
	clientCSR, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		log.Fatalf("parse CSR error: %v", err)
	}

	jsonData, err := json.MarshalIndent(clientCSR, "", "    ")
	if err != nil {
		log.Fatalf("JSON error: %v", err)
	}

	fmt.Printf("%s", jsonData)
}
