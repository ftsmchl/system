#!/bin/bash

./sysd >> sysd.out & -D status = $?
if [$status -ne 0]; then
	echo "Failed to start sysd process : $status"
fi

node renter_server/renter_server.js >> renter_server.out & -D status = $?
if [$status -ne 0]; then 
	echo "Failed to start renter_server : $status"
fi


#./sysclient accountAdd 0x6B054053c741495D0d42d9269E2e5C9dA985b969 >> client.out
#sleep 2
#./sysclient renter createContracts >> client.out

while sleep 60; do
	ps aux | grep sysd | grep -q -v grep
	STATUS_1=$?
	ps aux | grep node | grep -q -v grep
	STATUS_2=$?
	if [$STATUS_1 -ne 0 -o $STATUS_2 -ne 0]
		echo "One of the processes has already exited."
	fi
done
