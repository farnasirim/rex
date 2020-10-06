#!/bin/bash
openssl s_server -accept 7569 \
-CAfile ca.crt \
-cert $1 \
-key $2 \
-Verify 10 -tls1_3 -state -quiet
