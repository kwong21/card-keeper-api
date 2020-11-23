#!/bin/bash

COCKROACH_CLI="/cockroach/cockroach.sh"
HOSTPARAMS="start-single-node --insecure"

${COCKROACH_CLI} start-single-node --insecure --background
${COCKROACH_CLI} sql --insecure -e "CREATE USER IF NOT EXISTS tester";
${COCKROACH_CLI} sql --insecure -e "CREATE DATABASE CardKeeper;"
${COCKROACH_CLI} sql --insecure -e "GRANT ALL ON DATABASE CardKeeper to tester";

tail -f /dev/null