#!/bin/bash

./sysd >> logs/sysd.out &
#node renter_server/renter_server.js >> renter_server.out &

node host_server/host_server.js >> logs/host_server.out &

#./sysclient accountAdd 0x6B054053c741495D0d42d9269E2e5C9dA985b969 >> client.out
sleep 5
#./sysclient renter createContracts >> client.out


while sleep 60; do
	ps aux | grep sysd | grep -q -v grep
	STATUS_1=$?
	ps aux | grep node | grep -q -v grep
	STATUS_2=$?
	if [$STATUS_1 -ne 0 -o $STATUS_2 -ne 0]; then
		exit 1
	fi
done
