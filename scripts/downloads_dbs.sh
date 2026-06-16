#!/bin/bash

#
# Download DBs from https://github.com/sapics/ip-location-db/releases
#

set -e

# Create a directory to hold the databases if it doesn't exist
mkdir -p ./data

echo "*" > ./data/.gitignore

echo "Downloading DB-IP Country database..."
curl -L -o ./data/dbip-country.mmdb "https://github.com/sapics/ip-location-db/releases/download/latest/dbip-country-ipv4.mmdb"

echo "Downloading DB-IP ASN database..."
curl -L -o ./data/dbip-asn.mmdb "https://github.com/sapics/ip-location-db/releases/download/latest/dbip-asn-ipv4.mmdb"

echo "Downloading DB-IP City database..."
curl -L -o ./data/dbip-city.mmdb "https://github.com/sapics/ip-location-db/releases/download/latest/dbip-city-ipv4.mmdb"

echo "Databases successfully downloaded to ./data/"
