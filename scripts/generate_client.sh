#!/bin/bash

CERTS_DIR=certs/client
mkdir -p "$CERTS_DIR"

client_uuid=$(uuid)

openssl genrsa -out "$CERTS_DIR/$client_uuid.key" 2048
openssl req -new -key "$CERTS_DIR/$client_uuid.key" -out $client_uuid.csr -subj "/C=CA/ST=ON/O=UofT/OU=CS/CN=$client_uuid"
openssl x509 -req -in $client_uuid.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out "$CERTS_DIR/$client_uuid.pem" -days 30 -extfile openssl.cnf -extensions client_cert
