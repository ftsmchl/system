#!/bin/bash
echo "Destroying system's containers"

docker rm system_renter
docker rm system_host_1
docker rm system_host_2
