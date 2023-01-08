package beacon

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func loadCertificateFromFile(certFile string) (*x509.Certificate, error) {
	certEncoded, err := os.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("can not read certificate: %w", err)

	}

	certDecoded, _ := pem.Decode(certEncoded)
	certificate, err := x509.ParseCertificate(certDecoded.Bytes)
	if err != nil {
		return nil, fmt.Errorf("can not parse certificate, %w", err)
	}

	return certificate, nil
}

func GetTLSConfig(certFile string) (*tls.Config, error) {
	cp := x509.NewCertPool()

	certificate, err := loadCertificateFromFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("can not load credentials, %w", err)
	}

	cp.AddCert(certificate)

	return &tls.Config{
		RootCAs: cp,
	}, nil
}
