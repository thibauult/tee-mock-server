package pki

import (
	"crypto/rsa"
	"crypto/x509"
	"embed"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
)

var (
	//go:embed *.pem
	fs embed.FS
)

func GetSigningPrivateKey() (*rsa.PrivateKey, error) {
	keyData, err := fs.ReadFile("leaf-key.pem")
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}

	return privateKey, nil
}

func GetRootCertificate() string {
	certData, err := fs.ReadFile("root.pem")
	if err != nil {
		log.Fatalf("error reading root.pem file: %v", err)
	}
	return string(certData)
}

func GetCertificateChain() []string {
	root, err := getBase64EncodedDER("root.pem")
	if err != nil {
		log.Fatal(err)
	}
	intermediate, err := getBase64EncodedDER("intermediate.pem")
	if err != nil {
		log.Fatal(err)
	}
	leaf, err := getBase64EncodedDER("leaf.pem")
	if err != nil {
		log.Fatal(err)
	}
	return []string{leaf, intermediate, root}
}

func getBase64EncodedDER(name string) (string, error) {
	certData, err := fs.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("error reading private key file: %w", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(certData)
	if block == nil || block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("failed to decode PEM block containing the certificate")
	}

	// Parse the X.509 certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("error parsing X.509 certificate: %w", err)
	}

	// Get the Base64-encoded DER version of the certificate
	return base64.StdEncoding.EncodeToString(cert.Raw), nil
}
