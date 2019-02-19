#!/bin/bash -x

if ! [[ "$0" =~ "./gencerts.sh" ]]; then
	echo "must be run from 'fixtures'"
	exit 255
fi

if ! which cfssl; then
	echo "cfssl is not installed"
	exit 255
fi

SUFFIX="${SUFFIX-}"

cfssl gencert --initca=true ./ca-csr$SUFFIX.json | cfssljson --bare ./ca
mv ca.pem ca$SUFFIX.crt
openssl x509 -in ca$SUFFIX.crt -noout -text

# generate wildcard certificates DNS: *.etcd.local
cfssl gencert \
    --ca ./ca$SUFFIX.crt \
    --ca-key ./ca-key.pem \
    --config ./gencert.json \
    ./server-ca-csr$SUFFIX.json | cfssljson --bare ./server
mv server.pem server$SUFFIX.crt
mv server-key.pem server$SUFFIX.key.insecure

rm -f *.csr *.pem *.stderr *.txt

if [ -z "$WITH_SUFFIX" ]; then
	export WITH_SUFFIX=1
	export SUFFIX=-metrics
	"$0" "$@"
fi
