#!/bin/sh

sleep 3

/hati/build/bin/hati start --tcp --tcp-port 4242 --tcp-host 0.0.0.0 --rpc --rpc-port 6767 --rpc-host 0.0.0.0 --data-dir /etc/hati/data