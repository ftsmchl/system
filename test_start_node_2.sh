#!/bin/bash

./sysd >> logs/sysd.out &
sleep 1
case $HOSTNAME in
	(system_renter) 
		node renter_server/renter_server.js >> logs/renter_server.out &		
		sleep 1
		./sysclient accountAdd 0x11C593FF9D7fa0741645Cf60Ff0A67cA36BB9d2b >> logs/sysclient.out
		sleep 9 
		./sysclient renter createContracts >> logs/sysclient.out
		;;
	
	(system_host_1) 
		node host_server/host_server.js >> logs/host_server.out &	
		sleep 1
		./sysclient accountAdd  0xD7A7F523A228950eA91a328aC6d1AbF01e985802 >> logs/sysclient.out
		sleep 13 
		./sysclient host findContracts
		;;
	(system_host_2)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 1
		./sysclient accountAdd  0x1BF0ac46a277FA6cbaEe26512222334228440534 >> logs/sysclient.out
		sleep 14
		./sysclient host findContracts
		;;
	(system_host_3)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 1
		./sysclient accountAdd  0x1167E22ee3663e487975863556c49b2147cE4CeD >> logs/sysclient.out
		sleep 15 
		./sysclient host findContracts
		;;
	(system_host_4)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 1
		./sysclient accountAdd  0xEfE10d187A1588508558F0760D2f8D3aE8be894D >> logs/sysclient.out
		sleep 16 
		./sysclient host findContracts
		;;
	(system_host_5)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 1
		./sysclient accountAdd  0x6675AeD2016441b4410b3435F7742D970b3F35E8 >> logs/sysclient.out
		sleep 17 
		./sysclient host findContracts
		;;

esac

while sleep 60;do
	ps aux | grep sysd | grep -q -v grep
	STATUS_1=$?
	ps auz | grep node | grep -q -v grep
	STATUS_2=$?
	if [$STATUS_1 -ne 0 -o $STATUS_2 -ne 0]; then
		exit 1
	fi
done

