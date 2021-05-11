#!/bin/bash

./sysd >> logs/sysd.out &
sleep 1

case $HOSTNAME in
	(system_host_2)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 2
		./sysclient accountAdd  0x1BF0ac46a277FA6cbaEe26512222334228440534 >> logs/sysclient.out &
		sleep 6
		./sysclient host register >> logs/host_server.out &
		;;
	
	(system_host_3)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 9
		./sysclient accountAdd  0x1167E22ee3663e487975863556c49b2147cE4CeD >> logs/sysclient.out &
		sleep 13
		./sysclient host register >> logs/host_server.out &
		;;
	
	
	(system_host_4)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 17 
		./sysclient accountAdd  0xEfE10d187A1588508558F0760D2f8D3aE8be894D >> logs/sysclient.out &
		sleep 21 
		./sysclient host register >> logs/host_server.out &
		;;

	(system_host_5)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 25 
		./sysclient accountAdd  0x6675AeD2016441b4410b3435F7742D970b3F35E8 >> logs/sysclient.out & 
		sleep 29 
		./sysclient host register >> logs/host_server.out &
		;;

	(system_host_6)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 33 
		./sysclient accountAdd  0x1335E3c11621EE8c6cbe5d11Ab7910c1Ef62Cc53 >> logs/sysclient.out & 
		sleep 37 
		./sysclient host register >> logs/host_server.out &
		;;

esac

while sleep 60; do 
	ps aux | grep sysd | grep -q -v grep
	STATUS_1=$?
	ps aux | grep node | grep -q -v grep
	STATUS_2=$
	if [$STATUS_1 -ne 0 -o $STATUS_2 -ne 0]; then
		exit 1
	fi
done
