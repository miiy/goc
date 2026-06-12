package credentials

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	grpccredentials "google.golang.org/grpc/credentials"
)

func NewClientMTLS(serverName, certFilePath, keyFilePath, caFilePath string) (grpccredentials.TransportCredentials, error) {
	tlsConfig, err := ClientMTLSConfig(serverName, certFilePath, keyFilePath, caFilePath)
	if err != nil {
		return nil, err
	}
	return grpccredentials.NewTLS(tlsConfig), nil
}

func NewServerMTLS(certFilePath, keyFilePath, caFilePath string) (grpccredentials.TransportCredentials, error) {
	tlsConfig, err := ServerMTLSConfig(certFilePath, keyFilePath, caFilePath)
	if err != nil {
		return nil, err
	}
	return grpccredentials.NewTLS(tlsConfig), nil
}

func ClientMTLSConfig(serverName, certFilePath, keyFilePath, caFilePath string) (*tls.Config, error) {
	cert, ca, err := loadMTLSMaterials(certFilePath, keyFilePath, caFilePath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		ServerName:   serverName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}, nil
}

func ServerMTLSConfig(certFilePath, keyFilePath, caFilePath string) (*tls.Config, error) {
	cert, ca, err := loadMTLSMaterials(certFilePath, keyFilePath, caFilePath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
	}, nil
}

func loadMTLSMaterials(certFilePath, keyFilePath, caFilePath string) (tls.Certificate, *x509.CertPool, error) {
	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		return tls.Certificate{}, nil, fmt.Errorf("load key pair: %w", err)
	}

	ca := x509.NewCertPool()
	caBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		return tls.Certificate{}, nil, fmt.Errorf("read ca cert %q: %w", caFilePath, err)
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		return tls.Certificate{}, nil, fmt.Errorf("parse ca cert %q", caFilePath)
	}

	return cert, ca, nil
}
