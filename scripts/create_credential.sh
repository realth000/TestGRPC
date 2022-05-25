#!/bin/bash

if [ ! -f ca.conf -o ! -f server.conf ];then
	echo 'ca.conf or server.conf not found, exit'
fi

function auto_remove() {
        # [ $# -lt 1 ] && exit
        for arg in $@:
        do
                [ -f $arg ] && rm -rf $arg 1>/dev/null
        done
}

echo -n 'Remove old credentials...'
auto_remove ca.key ca.csr ca.crt ca.pem
auto_remove server.key server.csr server.pem
auto_remove client.key client.csr client.pem
echo ' OK'

echo -n 'Create CA credentials...'
openssl genrsa -out ca.key 4096 1>/dev/null 2>&1

openssl req \
-new \
-sha256 \
-out ca.csr \
-key ca.key \
-config ca.conf 1>/dev/null 2>&1 << EOF





EOF

openssl x509 \
-req \
-days 36500 \
-in ca.csr \
-signkey ca.key \
-out ca.crt \
1>/dev/null 2>&1

openssl req -new -x509 -key ca.key -out ca.pem -days 36500 1>/dev/null 2>&1 << EOF








EOF

echo ' OK'

echo -n 'Create server credentials...'
openssl genrsa -out server.key 4096 1>/dev/null 2>&1

openssl req \
-new \
-sha256 \
-out server.csr \
-key server.key \
-config server.conf 1>/dev/null 2>&1 << EOF






EOF

openssl x509 \
-req \
-days 365 \
-CA ca.crt \
-CAkey ca.key \
-CAcreateserial \
-in server.csr \
-out server.pem \
-extensions req_ext \
-extfile server.conf 1>/dev/null 2>&1

echo ' OK'

echo -n 'Create client credentials...'
openssl ecparam -genkey -name secp384r1 -out client.key 1>/dev/null 2>&1

openssl req -new -key client.key -out client.csr -config server.conf 1>/dev/null 2>&1 << EOF







EOF

openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem -extensions req_ext -extfile server.conf 1>/dev/null 2>&1
echo ' OK'
echo 'Success'
