#!/bin/sh

set -e

while true
do
    for i in "" / /v1/ /v1/Systems /v1/Systems/437XR1138R2 /v1/Systems/dummy
    do
        curl http://localhost:8080/redfish${i}   
    done
done