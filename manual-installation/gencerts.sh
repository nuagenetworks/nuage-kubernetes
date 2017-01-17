#!/bin/bash
# Parse CLI options
for i in "$@"; do
    case $i in
        --output-cert-dir=*)
            OUTDIR="${i#*=}"
            REST_SERVER_KEY=${OUTDIR}/nuageMonServer.key
            REST_SERVER_CERT=${OUTDIR}/nuageMonServer.crt
            REST_CLIENT_KEY=${OUTDIR}/nuageMonClient.key
            REST_CLIENT_CERT=${OUTDIR}/nuageMonClient.crt
            REST_CA_KEY=${OUTDIR}/nuageMonCA.key
            REST_CA_CERT=${OUTDIR}/nuageMonCA.crt
            REST_CA_SERIAL=${OUTDIR}/nuageMonCA.serial.txt
        ;;
    esac
done

# If any are missing, print the usage and exit
if [ -z $OUTDIR ]; then
    echo "Invalid syntax: $@"
    echo "Usage:"
    echo "  $0 --output-cert-dir=/path/to/output/dir/"
    echo "--output-cert-dir:  Directory to put artifacts in"
    echo ""
    echo "All options are required"
    exit 1
fi

openssl genrsa -out $REST_CA_KEY 4096
openssl req -new -x509 -key $REST_CA_KEY -out $REST_CA_CERT -subj "/CN=nuage-signer"
echo '00' > $REST_CA_SERIAL

OPENSSL_CONFIG='
[ clientauth ]
basicConstraints=CA:FALSE
extendedKeyUsage=critical,clientAuth
'

echo "$OPENSSL_CONFIG" > ${OUTDIR}/openssl.cnf
openssl genrsa -out $REST_SERVER_KEY 4096
openssl req -key $REST_SERVER_KEY -new -out ${OUTDIR}/restServer.req -subj "/CN=$(hostname -f)"
openssl x509 -req -in ${OUTDIR}/restServer.req -CA $REST_CA_CERT -CAkey \
    $REST_CA_KEY -CAserial $REST_CA_SERIAL \
    -out $REST_SERVER_CERT

openssl genrsa -out $REST_CLIENT_KEY 4096
openssl req -key $REST_CLIENT_KEY -new -out ${OUTDIR}/restClient.req -subj "/CN=nuage-client"
openssl x509 -req -in ${OUTDIR}/restClient.req -CA $REST_CA_CERT -CAkey \
    $REST_CA_KEY -CAserial $REST_CA_SERIAL \
    -out $REST_CLIENT_CERT -extensions clientauth \
    -extfile ${OUTDIR}/openssl.cnf

rm ${OUTDIR}/{restServer,restClient}.req

