#!/bin/bash
echo "Destroying system's containers"

docker rm system_renter
docker rm system_host_1
docker rm system_host_2
docker rm system_host_3
docker rm system_host_4
docker rm system_host_5
docker rm system_host_6
docker rm system_host_7
docker rm system_host_8
