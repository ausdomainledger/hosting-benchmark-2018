#!/usr/bin/env bash

BENCHMARKS=(describe_env bench micro_bench micro_cpu micro_io database)
BASEURL="$1"

if [ -z "$BASEURL" ]; then
  echo "Must pass a URL"
  exit 1
fi

for BENCHMARK in ${BENCHMARKS[*]}
do
  URL="$BASEURL$BENCHMARK.php"
  read -n1 -r -p "Press any key to run $URL live"
  curl "$URL" | tee "$BENCHMARK.results"
  read -n1 -r -p "Press any key to continue to next benchmark"
done