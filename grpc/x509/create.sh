#!/bin/bash

# Create the CA certs.
openssl req -x509                                     \
  -newkey rsa:4096                                    \
  -nodes                                              \
  -days 3650                                          \
  -keyout ca_key.pem                                  \
  -out ca_cert.pem                                    \
  -subj /C=US/ST=CA/L=SVL/O=gRPC/CN=goc-server_ca/    \
  -config ./openssl.cnf                               \
  -extensions goc_ca                                  \
  -sha256

# Generate a server cert.
openssl genrsa -out server_key.pem 4096
openssl req -new                                    \
  -key server_key.pem                               \
  -days 3650                                        \
  -out server_csr.pem                               \
  -subj /C=US/ST=CA/L=SVL/O=gRPC/CN=goc-server1/    \
  -config ./openssl.cnf                             \
  -reqexts goc_server
openssl x509 -req           \
  -in server_csr.pem        \
  -CAkey ca_key.pem         \
  -CA ca_cert.pem           \
  -days 3650                \
  -set_serial 1000          \
  -out server_cert.pem      \
  -extfile ./openssl.cnf    \
  -extensions goc_server    \
  -sha256
openssl verify -verbose -CAfile ca_cert.pem  server_cert.pem

# Generate a client cert.
openssl genrsa -out client_key.pem 4096
openssl req -new                                    \
  -key client_key.pem                               \
  -days 3650                                        \
  -out client_csr.pem                               \
  -subj /C=US/ST=CA/L=SVL/O=gRPC/CN=goc-client1/    \
  -config ./openssl.cnf                             \
  -reqexts goc_client
openssl x509 -req           \
  -in client_csr.pem        \
  -CAkey ca_key.pem         \
  -CA ca_cert.pem           \
  -days 3650                \
  -set_serial 1000          \
  -out client_cert.pem      \
  -extensions goc_client    \
  -sha256
openssl verify -verbose -CAfile ca_cert.pem  client_cert.pem

rm *_csr.pem
