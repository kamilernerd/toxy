#! /bin/bash

echo "Creating cerficates"

if openssl 2>/dev/null; then
	echo "Found openssl"
	openssl req -x509 -nodes -newkey rsa:2048 -keyout test/server.key -out test/server.crt -days 3650
	openssl req -new -sha256 -key test/server.key -out test/server.csr
	openssl x509 -req -sha256 -in test/server.csr -signkey test/server.key -out test/server.crt -days 3650
else
	echo "Openssl is missing please install."
	exit 1
fi
