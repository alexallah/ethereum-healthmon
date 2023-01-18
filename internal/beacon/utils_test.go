package beacon

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"testing"
	"time"
)

// generate a certificate and a private key
func genCertificate() ([]byte, []byte) {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"ORG"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}
	cert := &bytes.Buffer{}
	pem.Encode(cert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyBuf := &bytes.Buffer{}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	if err := pem.Encode(keyBuf, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Fatalf("Failed to write data to key.pem: %v", err)
	}

	return cert.Bytes(), keyBuf.Bytes()
}

func Test_GetTLSConfig(t *testing.T) {
	// save certificate
	dir := t.TempDir()

	cert, privKey := genCertificate()

	otherfile := dir + "/nocert.crt"
	err := os.WriteFile(otherfile, []byte("something else"), 0644)
	if err != nil {
		t.Error(err)
	}

	certfile := dir + "/cert.crt"
	err = os.WriteFile(certfile, cert, 0644)
	if err != nil {
		t.Error(err)
	}

	privkeyfile := dir + "/private.key"
	err = os.WriteFile(privkeyfile, privKey, 0644)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		file   string
		result bool
	}{
		{dir + "/nofilehere", false},
		{dir, false},
		{otherfile, false},
		{privkeyfile, false},
		{certfile, true},
	}

	for _, test := range tests {
		_, err := GetTLSConfig(test.file)
		if err != nil && test.result {
			t.Error("cert is not supposed to be loaded", test.file)
		}
		if err == nil && !test.result {
			t.Error("cert supposed to be loaded", test.file)
		}
	}
}
