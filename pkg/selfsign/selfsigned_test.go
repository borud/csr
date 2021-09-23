package selfsign

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSelfSigned(t *testing.T) {
	t.Parallel()

	ss, err := GenerateSelfSigned(Request{
		Hosts:        "127.0.0.1,example.com",
		Organization: "Test Inc",
		ValidFrom:    time.Now(),
		ValidUntil:   time.Now().Add(24 * time.Hour),
		IsCA:         true,
	})
	assert.Nil(t, err)
	assert.NotNil(t, ss)

	// Make sure we can load the cert and key
	cert, err := tls.X509KeyPair(ss.CertPEM, ss.KeyPEM)
	assert.Nil(t, err)
	assert.Greater(t, len(cert.Certificate), 0)
}
