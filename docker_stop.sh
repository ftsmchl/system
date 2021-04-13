#!/bin/bash

docker stop system_renter

for i in {1..8}
do
	docker stop system_host_$i
done
