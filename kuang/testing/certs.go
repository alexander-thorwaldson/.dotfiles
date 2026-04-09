// Package testing provides test helpers for kuang integration tests.
package testing

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"testing"
	"time"
)

// TestCA holds a self-signed CA and helper methods for issuing test certificates.
type TestCA struct {
	Cert    *x509.Certificate
	Pool    *x509.CertPool
	Key     *ecdsa.PrivateKey
	CertPEM []byte
}

// TestCert holds a certificate issued by a TestCA.
type TestCert struct {
	TLS     tls.Certificate
	Cert    *x509.Certificate
	Key     *ecdsa.PrivateKey
	CertPEM []byte
	KeyPEM  []byte
}

// NewTestCA generates a self-signed CA for testing.
func NewTestCA(t *testing.T) *TestCA {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating CA key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test-ca", Organization: []string{"kuang-test"}},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:         true,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("creating CA cert: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("parsing CA cert: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	pool := x509.NewCertPool()
	pool.AddCert(cert)

	return &TestCA{Cert: cert, Key: key, CertPEM: certPEM, Pool: pool}
}

// Issue creates a certificate signed by this CA.
func (ca *TestCA) Issue(t *testing.T, cn string, ou []string) *TestCert {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating key for %s: %v", cn, err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName:         cn,
			Organization:       []string{"kuang-test"},
			OrganizationalUnit: ou,
		},
		NotBefore:   time.Now().Add(-time.Hour),
		NotAfter:    time.Now().Add(24 * time.Hour),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:    []string{"localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, ca.Cert, &key.PublicKey, ca.Key)
	if err != nil {
		t.Fatalf("creating cert for %s: %v", cn, err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("parsing cert for %s: %v", cn, err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatalf("marshalling key for %s: %v", cn, err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("creating TLS cert for %s: %v", cn, err)
	}

	return &TestCert{Cert: cert, Key: key, CertPEM: certPEM, KeyPEM: keyPEM, TLS: tlsCert}
}

// MTLSClient creates an HTTP client that presents the given cert and trusts the CA.
func MTLSClient(ca *TestCA, cert *TestCert) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert.TLS},
				RootCAs:      ca.Pool,
				MinVersion:   tls.VersionTLS13,
			},
		},
	}
}
