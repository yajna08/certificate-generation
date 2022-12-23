package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

type Cert struct {
	certType  string
	certPrefx string
	certPath  string
}

type CertificateConfig struct {
	organization       []string
	country            []string
	province           []string
	locality           []string
	streetAddress      []string
	postalCode         []string
	Pincode            []string
	commonName         string
	organizationalUnit []string
	notBefore          time.Time
	notAfter           time.Time
	parent             string
	isCA               bool
	basePath           string
	certName           string
	payload            *bytes.Buffer
}

// Step 1
func newCertificateConfig(cm map[string]interface{}, crt *Cert) *CertificateConfig {

	var isCA bool
	var commonName string

	if crt.certType == "ca" || crt.certType == "tlsca" {
		isCA = true
	} else {
		isCA = false
	}
	commonName = crt.certPrefx + fmt.Sprintf("%v", cm["domain"])

	return &CertificateConfig{
		organization:       []string{fmt.Sprintf("%v", cm["organization"])},
		country:            []string{fmt.Sprintf("%v", cm["country"])},
		province:           []string{fmt.Sprintf("%v", cm["province"])},
		locality:           []string{fmt.Sprintf("%v", cm["locality"])},
		streetAddress:      []string{fmt.Sprintf("%v", cm["streetAddress"])},
		postalCode:         []string{fmt.Sprintf("%v", cm["postalCode"])},
		commonName:         commonName,
		organizationalUnit: []string{fmt.Sprintf("%v", cm["organizationalUnit"])},
		notBefore:          time.Now(),
		notAfter:           time.Now().AddDate(10, 0, 0),
		isCA:               isCA,
	}
}

// Step 2
func createX509Certificate(config *CertificateConfig) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:       config.organization,
			Country:            config.country,
			Province:           config.province,
			Locality:           config.locality,
			StreetAddress:      config.streetAddress,
			PostalCode:         config.postalCode,
			CommonName:         config.commonName,
			OrganizationalUnit: config.organizationalUnit,
		},
		NotBefore:             config.notBefore,
		NotAfter:              config.notAfter,
		IsCA:                  config.isCA,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

// Step 3
func generateECParamKeyCert(parent *x509.Certificate, child *x509.Certificate) (*bytes.Buffer, *bytes.Buffer) {
	keyCurve := elliptic.P256()

	//privateKey := new(ecdsa.PrivateKey)
	privateKey, err := ecdsa.GenerateKey(keyCurve, rand.Reader)
	publicKey := privateKey.PublicKey

	if err != nil {
		fmt.Println("error generating keys")
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, parent, child, &publicKey, privateKey)
	if err != nil {
		fmt.Println("error creating certificate, ", err)
	}

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	return caPrivKeyPEM, caPEM
}

// Step 4
func writePemAndKey(crt *Cert, cert *bytes.Buffer, key *bytes.Buffer) {

	// Saving Certificate
	if !(writeToFile(crt.certPath, cert)) {
		fmt.Println("unable to save certificate")
	}

	// Saving Key
	if !(writeToFile(crt.certPath, key)) {
		fmt.Println("unable to save key")
	}
}

// Step 5
func writeToFile(filepath string, content *bytes.Buffer) bool {
	// Write a string to a file
	err := os.WriteFile(filepath, content.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func main() {

	// These variables will come dynamically from network-config.json file
	peerCount := 2
	userCount := 2

	// Common Cert Config
	commonParams := map[string]interface{}{
		"organization":       "Org1",
		"country":            "India",
		"province":           "Karnataka",
		"locality":           "Bengaluru",
		"streetAddress":      "Murgespalya Road",
		"postalCode":         "560071",
		"organizationalUnit": "Org1",
		"domain":             "org1.example.com",
	}
	// Generting Root CA Config
	rootCAConfig := &Cert{"ca", "ca", "./"}
	rootCACert := createX509Certificate(newCertificateConfig(commonParams, rootCAConfig))
	rootCAPem, rootCAKey := generateECParamKeyCert(rootCACert, rootCACert)
	writePemAndKey(rootCAConfig, rootCAPem, rootCAKey)

	// Generating TLS CA Config
	tlsCAConfig := &Cert{"tlsca", "tlsca", "./"}
	tlsCACert := createX509Certificate(newCertificateConfig(commonParams, tlsCAConfig))
	tlsCAPem, tlsCAKey := generateECParamKeyCert(tlsCACert, tlsCACert)
	writePemAndKey(tlsCAConfig, tlsCAPem, tlsCAKey)

	// Generating Admin Config
	adminConfig := &Cert{"admin", "admin", "./"}
	adminCert := createX509Certificate(newCertificateConfig(commonParams, adminConfig))
	adminCAPem, adminCAKey := generateECParamKeyCert(rootCACert, adminCert)
	writePemAndKey(adminConfig, adminCAPem, adminCAKey)
	admintlsCAPem, admintlsCAKey := generateECParamKeyCert(tlsCACert, adminCert)
	writePemAndKey(adminConfig, admintlsCAPem, admintlsCAKey)

	// Generating Peer Config
	for i := 0; i < peerCount; i++ {
		peerConfig := &Cert{"peer", "peer" + fmt.Sprintf("%v", i), "./"}
		peerCACert := createX509Certificate(newCertificateConfig(commonParams, peerConfig))
		peerCAPem, peerCAKey := generateECParamKeyCert(rootCACert, peerCACert)
		writePemAndKey(peerConfig, peerCAPem, peerCAKey)
		peerTLSCAPem, peerTLSCAKey := generateECParamKeyCert(tlsCACert, peerCACert)
		writePemAndKey(peerConfig, peerTLSCAPem, peerTLSCAKey)
	}

	// Generating User Config
	for i := 0; i < userCount; i++ {
		userConfig := &Cert{"user", "user" + fmt.Sprintf("%v", i), "./"}
		userCACert := createX509Certificate(newCertificateConfig(commonParams, userConfig))
		userCAPem, userCAKey := generateECParamKeyCert(rootCACert, userCACert)
		writePemAndKey(userConfig, userCAPem, userCAKey)
		userTLSCAPem, userTLSCAKey := generateECParamKeyCert(tlsCACert, userCACert)
		writePemAndKey(userConfig, userTLSCAPem, userTLSCAKey)
	}
}
