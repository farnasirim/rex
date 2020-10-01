## Scripts

### Generating certificates
First, run `./create_ca.sh` to create a CA in the current directory (`scripts`).
Then run `./generate_server.sh` to create and sign server's key pair and
`./generate_client.sh` as many times as required to create and client
certificates. Files go to `./certs/{client,server}`.

You can test drive the generated certificates using openssl to create a secure
channel. Run `./test_server.sh path/to/cert path/to/key` and
`./test_client.sh path/to/cert path/to/key`.
