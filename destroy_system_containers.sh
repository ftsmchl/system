#!/bin/bash
echo "Destroying system's containers"

docker rm system_renter

for i in {1..8}
do
	docker rm system_host_$i
done
