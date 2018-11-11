#!/usr/bin/env bash

set -euf -o pipefail

# Download and build wrk2
if [ ! -d "./wrk2" ]; then
  git clone https://github.com/giltene/wrk2
fi
if [ ! -f "./wrk2/wrk" ]; then
  pushd ./wrk2
  make
  popd
fi

# Parameters
while getopts ":u:r:" o; do
  case "${o}" in
    u)
      URL=${OPTARG}
      ;;
    r)
      RATE=${OPTARG}
      ;;
  esac
done
if [ -z "${URL:=}" ] || [ -z "${RATE:=}" ]; then
    echo "Usage: $0 -u URL [http://example.org] -r RATE [1,5,10]"
    exit 1
fi

echo "Benchmarking ${URL} with ${RATE}/sec requests"

URL_SAFE=$(echo $URL | tr -cd '[[:alnum:]]._-')
./wrk2/wrk --latency -d60s -R${RATE} -t 1 -c ${RATE} ${URL} 2>&1 | tee "wordpress-results-${URL_SAFE}-${RATE}-$(date +'%Y%m%d%H%M%S%z').log"
