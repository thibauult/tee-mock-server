#!/bin/bash

# This script generates a valid X.509 certificate chain using OpenSSL, consisting of:
# - A Root CA certificate (self-signed)
# - An Intermediate CA certificate (signed by the Root CA)
# - A Leaf certificate (signed by the Intermediate CA)
#
# The following files will be created:
# - root-key.pem: Private key for the Root CA
# - root.pem: Root CA certificate
# - intermediate-key.pem: Private key for the Intermediate CA
# - intermediate.pem: Intermediate CA certificate
# - leaf-key.pem: Private key for the Leaf certificate
# - leaf.pem: Leaf certificate
#
# The script includes verification of the certificate chain and cleans up intermediate files (CSR files).
# Ensure that the file `intermediate-ext.cnf` exists in the same directory as this script with the necessary
# basicConstraints and other extensions for the Intermediate CA.
#
# Usage: Run this script directly in a bash environment.

# Define the intermediate-ext.cnf content as a constant
INTERMEDIATE_EXT="
[ v3_ca ]
basicConstraints = CA:TRUE
keyUsage = critical, keyCertSign, cRLSign
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer
"

# Generate the private key for the root CA
openssl genrsa -out root-key.pem 2048

# Create the root certificate
openssl req -x509 -new -nodes -key root-key.pem -sha256 -days 3650 -out root.pem -subj "/C=FR/ST=IDF/L=Paris/O=Symphony/CN=RootCA"

# Generate the private key for the intermediate CA
openssl genrsa -out intermediate-key.pem 2048

# Create the CSR for the intermediate CA
openssl req -new -key intermediate-key.pem -out intermediate.csr -subj "/C=FR/ST=IDF/L=Paris/O=Symphony/CN=IntermediateCA"

# Check if the temporary configuration file already exists
if [ ! -e intermediate-ext.cnf ]; then
  # Write the intermediate-ext.cnf content to a temporary file
  echo "Creating the temporary intermediate-ext.cnf file"
  echo "$INTERMEDIATE_EXT" > intermediate-ext.cnf
fi

# Sign the intermediate certificate with the root CA
openssl x509 -req -in intermediate.csr -CA root.pem -CAkey root-key.pem -CAcreateserial -out intermediate.pem -days 3650 -sha256 -extfile intermediate-ext.cnf -extensions v3_ca

# Generate the private key for the leaf certificate
openssl genrsa -out leaf-key.pem 2048

# Create the CSR for the leaf certificate
openssl req -new -key leaf-key.pem -out leaf.csr -subj "/C=FR/ST=IDF/L=Paris/O=Symphony/CN=LeafCertificate"

# Sign the leaf certificate with the intermediate CA
openssl x509 -req -in leaf.csr -CA intermediate.pem -CAkey intermediate-key.pem -CAcreateserial -out leaf.pem -days 3650 -sha256

# Verify the certificate chain
openssl verify -CAfile <(cat root.pem intermediate.pem) leaf.pem

# Cleanup
echo "Removing all .csr and .cnf files"
rm ./*.csr ./*.cnf
