#!/bin/bash

CERTS_DIR=certs/server
mkdir -p "$CERTS_DIR"

server_uuid=$(uuid)

openssl genrsa -out "$CERTS_DIR/$server_uuid.key" 2048
openssl req -new -key "$CERTS_DIR/$server_uuid.key" -out $server_uuid.csr -subj "/C=CA/ST=ON/O=UofT/OU=CS/CN=$server_uuid"
openssl x509 -req -in $server_uuid.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out "$CERTS_DIR/$server_uuid.pem" -days 30 -extfile openssl.cnf -extensions server_cert


