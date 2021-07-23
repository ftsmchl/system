#!/bin/bash

docker stop system_renter

pid=""

for i in {1..8}
do
	docker stop system_host_$i &

	pid="$pid $!"
done


for p in $pid
do
	wait $p
done
