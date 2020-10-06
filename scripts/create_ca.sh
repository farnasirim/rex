#!/bin/bash
# TODO: password protect

mkdir private newcerts
touch index.txt

openssl genrsa -out ca.key 4096
openssl req -config openssl.cnf -days 365 -new -x509 -key ca.key -out ca.crt -subj "/C=CA/ST=ON/O=UofT/OU=CS/CN=rexca"
