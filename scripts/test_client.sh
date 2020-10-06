#!/bin/bash
openssl s_client -connect localhost:7569 \
-CAfile ca.crt \
-cert $1 \
-key $2 \
-tls1_3 -state -quiet
