#!/bin/bash
# TODO: password protect

mkdir private newcerts
touch index.txt

openssl genrsa -out ca.key 4096
openssl req -config openssl.cnf -days 365 -new -x509 -key ca.key -out ca.crt -subj "/C=CA/ST=ON/O=UofT/OU=CS/CN=rexca"

server_uuid=$(uuid)

openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj "/C=CA/ST=ON/O=UofT/OU=CS/CN=$server_uuid"
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.pem -days 30 -extfile openssl.cnf -extensions server_cert


client_uuid=$(uuid)
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -subj "/C=CA/ST=ON/O=UofT/OU=CS/CN=$client_uuid"
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.pem -days 30 -extfile openssl.cnf -extensions client_cert
