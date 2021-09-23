package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}
var emailAddress = "user@example.com"

func main() {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("error generating key: %v", err)
	}

	// Create public key PEM
	publicBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Fatalf("error marshalling public key: %v", err)
	}
	publicPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	})
	fmt.Printf("%s\n", publicPem)

	// Create private key PEM
	privateBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("error marshalling private key: %v", err)
	}
	privatePem := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateBytes,
	})
	fmt.Printf("%s\n", privatePem)

	subject := pkix.Name{
		CommonName: "sample client certificate",
	}

	rawSubject := subject.ToRDNSequence()

	rawSubject = append(rawSubject, []pkix.AttributeTypeAndValue{
		{
			Type:  oidEmailAddress,
			Value: emailAddress,
		},
	})

	asn1Subj, err := asn1.Marshal(rawSubject)
	if err != nil {
		log.Fatalf("error marshalling subject to asn1: %v", err)
	}

	template := x509.CertificateRequest{
		RawSubjectPublicKeyInfo: pub,
		RawSubject:              asn1Subj,
		SignatureAlgorithm:      x509.PureEd25519,
		PublicKey:               pub,
		Subject:                 subject,
		EmailAddresses:          []string{emailAddress},
	}

	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, priv)
	csrPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	fmt.Printf("csrBytes size: %d\n", len(csrBytes))
	fmt.Printf("  csrPEM size: %d\n\n", len(csrPem))

	fmt.Printf("%s\n", csrPem)

	resp, err := http.Post("http://localhost:8881/sign", "application/x-pem-file", bytes.NewBuffer(csrPem))
	if err != nil {
		log.Fatalf("error performing POST to server: %v", err)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}
	resp.Body.Close()

	fmt.Printf("Client certificate signed by server:\n%s\n", responseBody)

	block, _ := pem.Decode(responseBody)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("error parsing certificate: %v", err)
	}

	fmt.Printf("PEM size: %d\n", len(responseBody))
	fmt.Printf("certificate size: %d\n", len(block.Bytes))
	fmt.Printf("Issuer: %s\n", cert.Issuer)
	fmt.Printf("Authority Key ID: %x\n", cert.AuthorityKeyId)
	fmt.Printf("Public key algorithm: %s\n", cert.PublicKeyAlgorithm)
}
