package main

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/borud/csr/pkg/selfsign"
	"github.com/gorilla/mux"
)

var (
	ca     *selfsign.SelfSigned
	caCert *x509.Certificate
)

const (
	listenAddr = ":8881"
)

func main() {
	var err error

	// Create CA certificate (temporary and self signed for demo purposes)
	ca, err = selfsign.GenerateSelfSigned(selfsign.Request{
		Organization: "Blind Faith Inc",
		Hosts:        "localhost,127.0.0.1",
		ValidFrom:    time.Now(),
		ValidUntil:   time.Now().Add(24 * 365 * time.Hour),
		IsCA:         true,
	})
	if err != nil {
		log.Fatalf("unable to generate self signed CA certificate: %v", err)
	}

	caCert, err = x509.ParseCertificate(ca.DERBytes)
	if err != nil {
		log.Fatalf("error parsing CA certificate: %v", err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/sign", SignHandler).Methods("POST")

	server := http.Server{
		Addr:    listenAddr,
		Handler: m,
	}

	log.Printf("server up, listening to %s", listenAddr)
	log.Print(server.ListenAndServe())
}

// SignHandler accepts a CSR and produces a signed certificate.
func SignHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}

	pemBlock, _ := pem.Decode(body)
	clientCSR, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		http.Error(w, "parse CSR failed", http.StatusBadRequest)
		log.Printf("parse CSR failed: %v", err)
	}
	log.Printf("Got CSR from %s with signature %x ", clientCSR.Subject.String(), clientCSR.Signature)

	err = clientCSR.CheckSignature()
	if err != nil {
		http.Error(w, "signature check failed", http.StatusBadRequest)
		log.Printf("signature check failed: %v", err)
	}
	log.Printf("signature ok")

	clientCRTTemplate := x509.Certificate{
		Signature:          clientCSR.Signature,
		SignatureAlgorithm: clientCSR.SignatureAlgorithm,
		PublicKeyAlgorithm: clientCSR.PublicKeyAlgorithm,
		PublicKey:          clientCSR.PublicKey,
		SerialNumber:       big.NewInt(2),
		Issuer:             caCert.Issuer,
		Subject:            clientCSR.Subject,
		NotBefore:          time.Now(),
		NotAfter:           time.Now().Add(24 * time.Hour),
		KeyUsage:           x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	clientCRTRaw, err := x509.CreateCertificate(rand.Reader, &clientCRTTemplate, caCert, clientCSR.PublicKey, ca.PrivateKey)
	if err != nil {
		http.Error(w, "error creating certificate", http.StatusBadRequest)
		log.Printf("error creating certificate: %v", err)
	}

	clientCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: clientCRTRaw,
	})

	log.Printf("created certificate:\n%s", clientCertPEM)

	w.Write(clientCertPEM)
}
