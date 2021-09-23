package selfsign

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"strings"
	"time"
)

// Request for self signed certificate.
type Request struct {
	Hosts        string
	Organization string
	ValidFrom    time.Time
	ValidUntil   time.Time
	IsCA         bool
}

// SelfSigned certificate and its key.
type SelfSigned struct {
	PrivateKey ed25519.PrivateKey
	DERBytes   []byte
	KeyPEM     []byte
	CertPEM    []byte
}

// GenerateSelfSigned creates a self signed certificate and returns certificate and key
func GenerateSelfSigned(c Request) (*SelfSigned, error) {
	if c.Hosts == "" {
		return nil, errors.New("no hostname given")
	}

	if c.ValidUntil.Before(c.ValidFrom) {
		return nil, fmt.Errorf("validUntil (%s) is before validFrom (%s)", c.ValidUntil.Format(time.RFC3339), c.ValidFrom.Format(time.RFC3339))
	}

	// Create serial number for X.509 certificate
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	// Create X.509 certificate template
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{c.Organization},
		},
		NotBefore:             c.ValidFrom,
		NotAfter:              c.ValidUntil,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  c.IsCA,
	}

	if c.IsCA {
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	// Add the hostnames and IPs
	hosts := strings.Split(c.Hosts, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
			continue
		}
		template.DNSNames = append(template.DNSNames, h)
	}

	// Create key and then certificate
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public().(ed25519.PublicKey), priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal private key: %w", err)
	}

	return &SelfSigned{
		PrivateKey: priv,
		DERBytes:   derBytes,
		CertPEM: pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derBytes,
		}),
		KeyPEM: pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privBytes,
		}),
	}, nil
}
