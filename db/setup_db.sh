#!/bin/bash
echo Wait for servers to be up
sleep 10

HOSTPARAMS="--host roach1 --insecure"
SQL="/cockroach/cockroach.sh sql $HOSTPARAMS"

$SQL < /cockroach/schema.sql
