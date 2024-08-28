#!/bin/bash -e

if [[ $1 = "-h" ]] || [[ $# -eq 0 ]] ; then
  echo "generate self-signed certificates for (m)TLS test"
  echo "    param1/param2 for rsa, output in 3/4"
  echo "    param3/param4 for ecdsa, output in 5/6"
  echo ""
  echo "Usage:"
  echo " script hostname password"
  echo " script hostname password"
  echo " script hostname password ecdsa-hostname ecdsa-passwd"
  echo " script clean       -- clean certs in 3/4/5/6"
  echo " script cleanall    -- clean cas in 1/2 and certs in 3/4/5/6"

  exit 0  
fi    
    
if [[ $1 = "clean" ]]; then
  rm -rf 3_application
  rm -rf 4_client
  rm -rf 5_application_ecdsa
  rm -rf 6_client_ecdsa
  
  exit 0
fi
if [[ $1 = "cleanall" ]]; then
  rm -rf 1_root
  rm -rf 2_intermediate
  rm -rf 3_application
  rm -rf 4_client
  rm -rf 5_application_ecdsa
  rm -rf 6_client_ecdsa
  
  exit 0
fi


if [ ! -d 1_root ]; then

echo 
echo Generate the root key
echo ---
mkdir -p 1_root/private
openssl genrsa -aes256 -passout pass:admin -out 1_root/private/ca.key.pem 4096

chmod 444 1_root/private/ca.key.pem


echo 
echo Generate the root certificate
echo ---
mkdir -p 1_root/certs
mkdir -p 1_root/newcerts
touch 1_root/index.txt
echo "100212" > 1_root/serial
openssl req -config openssl.cnf \
      -key 1_root/private/ca.key.pem \
      -passin pass:admin \
      -new -x509 -days 7300 -sha256 -extensions v3_ca \
      -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=test.com" \
      -out 1_root/certs/ca.cert.pem


echo 
echo Verify root key
echo ---
openssl x509 -noout -text -in 1_root/certs/ca.cert.pem

echo 
echo Generate the key for the intermediary certificate
echo ---
mkdir -p 2_intermediate/private
openssl genrsa -aes256 \
  -passout pass:admin \
  -out 2_intermediate/private/intermediate.key.pem 4096

chmod 444 2_intermediate/private/intermediate.key.pem


echo 
echo Generate the signing request for the intermediary certificate
echo ---
mkdir -p 2_intermediate/csr
openssl req -config openssl.cnf -new -sha256 \
      -passin pass:admin \
      -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=test.com" \
      -key 2_intermediate/private/intermediate.key.pem \
      -out 2_intermediate/csr/intermediate.csr.pem


echo 
echo Sign the intermediary
echo ---
mkdir -p 2_intermediate/certs
mkdir -p 2_intermediate/newcerts
touch 2_intermediate/index.txt
echo "100212" > 2_intermediate/serial
openssl ca -config openssl.cnf -extensions v3_intermediate_ca \
        -passin pass:admin \
        -days 36500 -notext -md sha256 \
        -in 2_intermediate/csr/intermediate.csr.pem \
        -out 2_intermediate/certs/intermediate.cert.pem

cat 2_intermediate/certs/intermediate.cert.pem 1_root/certs/ca.cert.pem > 2_intermediate/certs/ca-chain.cert.pem

chmod 444 2_intermediate/certs/intermediate.cert.pem
chmod 444 2_intermediate/certs/ca-chain.cert.pem


echo 
echo Verify intermediary
echo ---
openssl x509 -noout -text \
      -in 2_intermediate/certs/intermediate.cert.pem

openssl verify -CAfile 1_root/certs/ca.cert.pem \
      2_intermediate/certs/intermediate.cert.pem


echo 
echo Create the chain file
echo ---
cat 2_intermediate/certs/intermediate.cert.pem \
      1_root/certs/ca.cert.pem > 2_intermediate/certs/ca-chain.cert.pem

chmod 444 2_intermediate/certs/ca-chain.cert.pem

fi


if [ ! -d 3_application ]; then

echo 
echo Create the application key
echo ---
mkdir -p 3_application/private
openssl genrsa \
      -passout pass:$2 \
      -out 3_application/private/$1.key.pem 2048

chmod 444 3_application/private/$1.key.pem


echo 
echo Create the application signing request
echo ---
mkdir -p 3_application/csr
openssl req -batch -config intermediate_openssl.cnf \
      -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=$1" \
      -passin pass:$2 \
      -key 3_application/private/$1.key.pem \
      -new -sha256 -out 3_application/csr/$1.csr.pem


echo 
echo Create the application certificate
echo ---
mkdir -p 3_application/certs
openssl ca -batch -config intermediate_openssl.cnf \
      -passin pass:admin \
      -extensions server_cert -days 375 -notext -md sha256 \
      -in 3_application/csr/$1.csr.pem \
      -out 3_application/certs/$1.cert.pem

chmod 444 3_application/certs/$1.cert.pem


echo 
echo Validate the certificate
echo ---
openssl x509 -noout -text \
      -in 3_application/certs/$1.cert.pem


echo 
echo Validate the certificate has the correct chain of trust
echo ---
openssl verify -CAfile 2_intermediate/certs/ca-chain.cert.pem \
      3_application/certs/$1.cert.pem


echo
echo Generate the client key
echo ---
mkdir -p 4_client/private
openssl genrsa \
    -passout pass:$2 \
    -out 4_client/private/$1.key.pem 2048

chmod 444 4_client/private/$1.key.pem


echo
echo Generate the client signing request
echo ---
mkdir -p 4_client/csr
openssl req -batch -config intermediate_openssl.cnf \
      -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=$1" \
      -passin pass:$2 \
      -key 4_client/private/$1.key.pem \
      -new -sha256 -out 4_client/csr/$1.csr.pem


echo 
echo Create the client certificate
echo ---
mkdir -p 4_client/certs
openssl ca -batch -config intermediate_openssl.cnf \
      -passin pass:admin \
      -extensions usr_cert -days 375 -notext -md sha256 \
      -in 4_client/csr/$1.csr.pem \
      -out 4_client/certs/$1.cert.pem

chmod 444 4_client/certs/$1.cert.pem

fi


if [ $# -ge 4 ]; then
if [ ! -d 5_application_ecdsa ]; then

echo 
echo Create the application key
echo ---
mkdir -p 5_application_ecdsa/private
openssl ecparam -name P-256 -genkey \
    -out 5_application_ecdsa/private/$3.key.pem 

chmod 444 5_application_ecdsa/private/$3.key.pem


echo 
echo Create the application signing request
echo ---
mkdir -p 5_application_ecdsa/csr
openssl req -batch -config intermediate_openssl.cnf \
      -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=$3" \
      -passin pass:$4 \
      -key 5_application_ecdsa/private/$3.key.pem \
      -new -sha256 -out 5_application_ecdsa/csr/$3.csr.pem


echo 
echo Create the application certificate
echo ---
mkdir -p 5_application_ecdsa/certs
openssl ca -batch -config intermediate_openssl.cnf \
      -passin pass:admin \
      -extensions server_cert -days 375 -notext -md sha256 \
      -in 5_application_ecdsa/csr/$3.csr.pem \
      -out 5_application_ecdsa/certs/$3.cert.pem

chmod 444 5_application_ecdsa/certs/$3.cert.pem


echo 
echo Validate the certificate
echo ---
openssl x509 -noout -text \
      -in 5_application_ecdsa/certs/$3.cert.pem


echo 
echo Validate the certificate has the correct chain of trust
echo ---
openssl verify -CAfile 2_intermediate/certs/ca-chain.cert.pem \
      5_application_ecdsa/certs/$3.cert.pem


echo
echo Generate the client key
echo ---
mkdir -p 6_client_ecdsa/private
openssl ecparam -name P-256 -genkey \
    -out 6_client_ecdsa/private/$3.key.pem 

chmod 444 6_client_ecdsa/private/$3.key.pem


echo
echo Generate the client signing request
echo ---
mkdir -p 6_client_ecdsa/csr
openssl req -batch -config intermediate_openssl.cnf \
      -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=$3" \
      -passin pass:$4 \
      -key 6_client_ecdsa/private/$3.key.pem \
      -new -sha256 -out 6_client_ecdsa/csr/$3.csr.pem


echo 
echo Create the client certificate
echo ---
mkdir -p 6_client_ecdsa/certs
openssl ca -batch -config intermediate_openssl.cnf \
      -passin pass:admin \
      -extensions usr_cert -days 375 -notext -md sha256 \
      -in 6_client_ecdsa/csr/$3.csr.pem \
      -out 6_client_ecdsa/certs/$3.cert.pem

chmod 444 6_client_ecdsa/certs/$3.cert.pem

fi
fi


