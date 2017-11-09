package encryption

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

// CreateTLSConfig creates TLS configuration
func CreateTLSConfig(pemFile, pemCertFile, pemPrivateKeyFile string) *tls.Config {
	if len(pemFile) <= 0 {
		return nil
	}
	rootCerts, err := loadCertificates(pemFile)
	if err != nil {
		log.Fatalf("Couldn't load root certificate from %s. Got %s.", pemFile, err)
	}
	if len(pemCertFile) > 0 && len(pemPrivateKeyFile) > 0 {
		clientPrivateKey, err := loadPrivateKeyFrom(pemCertFile, pemPrivateKeyFile)
		if err != nil {
			log.Fatalf("Couldn't setup client authentication. Got %s.", err)
		}
		return &tls.Config{
			RootCAs:      rootCerts,
			Certificates: []tls.Certificate{*clientPrivateKey},
		}
	}
	return &tls.Config{
		RootCAs: rootCerts,
	}
}

func loadCertificates(pemFile string) (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}
	certificates := x509.NewCertPool()
	certificates.AppendCertsFromPEM(caCert)
	return certificates, nil
}

func loadPrivateKeyFrom(pemCertFile, pemPrivateKeyFile string) (*tls.Certificate, error) {
	privateKey, err := tls.LoadX509KeyPair(pemCertFile, pemPrivateKeyFile)
	if err != nil {
		return nil, err
	}
	return &privateKey, nil
}
