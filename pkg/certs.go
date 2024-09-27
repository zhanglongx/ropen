package pkg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

// Certs represents the CA certificate and private key,
// issuer of the website certificates.
// TODO: phase-key
type Certs struct {
	caCert *x509.Certificate
	caKey  *rsa.PrivateKey
}

func NewCerts(caPath, keyPath string) (*Certs, error) {
	// FIXME: PKCS1/12?
	caCertPEM, err := os.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	caBlock, _ := pem.Decode(caCertPEM)
	if caBlock == nil || caBlock.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	caCert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return nil, err
	}

	caKeyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	keyBlock, _ := pem.Decode(caKeyPEM)
	if keyBlock == nil || keyBlock.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to parse CA private key")
	}

	caKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// FIXME: cakey type
	return &Certs{
		caCert: caCert,
		caKey:  caKey.(*rsa.PrivateKey),
	}, nil
}

func (c *Certs) GenerateWebsiteCerts(ip string) (tls.Certificate, error) {
	serialNumber, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return tls.Certificate{}, err
	}

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{APP_NAME},
			OrganizationalUnit: []string{APP_NAME + " Unit"},
			CommonName:         ip,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		BasicConstraintsValid: true,
	}

	template.IPAddresses = []net.IP{net.ParseIP(ip)}

	certDER, err := x509.CreateCertificate(rand.Reader,
		&template, c.caCert, &privKey.PublicKey, c.caKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey)})

	return tls.X509KeyPair(certPEM, keyPEM)
}
