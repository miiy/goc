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

func ClientMTLSConfig(serverName, certFilePath, keyFilePath, caFilePath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("load client cert: %w", err)
	}

	ca := x509.NewCertPool()
	caBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		return nil, fmt.Errorf("read ca cert %q: %w", caFilePath, err)
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		return nil, fmt.Errorf("parse ca cert %q", caFilePath)
	}

	return &tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}, nil
}
