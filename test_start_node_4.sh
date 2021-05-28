#!/bin/bash

./sysd >> logs/sysd.out &
sleep 1
case $HOSTNAME in
	(system_renter) 
		node renter_server/renter_server.js >> logs/renter_server.out &		
		sleep 1
		./sysclient accountAdd  0x5aff1d56bC44898c28f94Ba6018cBc49718a666E 0xa82daf7b03ed8307357319d367015c242cdba49d4be1298ca4c2e330d397d59e >> logs/sysclient.out &
		sleep 9 
		./sysclient renter createContracts >> logs/sysclient.out
		sleep 65
		./sysclient renter uploadFile pame >> logs/sysclient.out 
		;;
	
	(system_host_1) 
		node host_server/host_server.js >> logs/host_server.out &	
		sleep 1
		./sysclient accountAdd  0xD7A7F523A228950eA91a328aC6d1AbF01e985802 0x2b34a6bae46d818bf316af17d293b43b81ac90a1a14c9de65faa400a52ab072c >> logs/sysclient.out &
		sleep 19 
		./sysclient host findContracts &
		#sleep 65
		#./sysclient host register >> logs/sysclient.out
		;;
	(system_host_2)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 1
		./sysclient accountAdd  0x1BF0ac46a277FA6cbaEe26512222334228440534 0x9c13227e6e38aa329a340a795d08810209de10313569aa00e2012bad21d17d25 >> logs/sysclient.out &
		sleep 24 
		./sysclient host findContracts &
		;;
	(system_host_3)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 1
		./sysclient accountAdd  0x1167E22ee3663e487975863556c49b2147cE4CeD 0x166db04561b19a9bddab4ffe8862c25e748e7f7080fb5deab3ab4a5367c27404 >> logs/sysclient.out &
		sleep 29 
		./sysclient host findContracts &
		;;
	(system_host_4)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 1
		./sysclient accountAdd  0xEfE10d187A1588508558F0760D2f8D3aE8be894D 0x3ded51d54c1009d865ab1f3caeba2a719e207113b51c2e02bd668c8654ed602d >> logs/sysclient.out &
		sleep 35 
		./sysclient host findContracts &
		;;
	(system_host_5)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 1
		./sysclient accountAdd  0x6675AeD2016441b4410b3435F7742D970b3F35E8 0x59a4f155e780b3b12fc275e17271fd987655133cb81971974018812dc8951b03 >> logs/sysclient.out &
		sleep 40 
		./sysclient host findContracts &
		;;
	(system_host_6)
		node host_server/host_server.js >> logs/host_server.out &
		sleep 1
		./sysclient accountAdd  0x1335E3c11621EE8c6cbe5d11Ab7910c1Ef62Cc53 0x1c31dce4ca0ce84d2efa2b98513269673d29229211021c70532c54403b716c18 >> logs/sysclient.out &
		sleep 49 
		./sysclient host findContracts &
		;;


esac

while sleep 60;do
	ps aux | grep sysd | grep -q -v grep
	STATUS_1=$?
	ps aux | grep node | grep -q -v grep
	STATUS_2=$?
	if [$STATUS_1 -ne 0 -o $STATUS_2 -ne 0]; then
		exit 1
	fi
done

