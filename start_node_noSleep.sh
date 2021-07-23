#!/bin/bash

./sysd >> logs/sysd.out &

cs=1

echo "counter1 = $cs" >> logs/counter1.out

while [ $cs -gt 0 ]
do
	echo "eimai sto while" >> logs/while.out
	result=$(nc -z -v localhost 8080 2>&1)
	if grep -q "succeeded" <<< "$result";
	then
		cs=$($cs - 1)
	fi

	echo $result >> logs/daem.out
done	


case $HOSTNAME in
	(system_renter) 
		node renter_server/renter_server.js >> logs/renter_server.out &		
		counter=1

		while [ $counter -gt 0 ]
		do
			result=$(nc -z -v localhost 8000 2>&1)
			if grep -q "succeeded" <<< "$result";
			then
				counter=$(($counter - 1))
			fi
		done

		#result=$(./sysclient accountAdd  0x8e169AbBf12dd68433E3bB1Fb0cd20Fd50263F11 0x773aee59f8099107c89d690e2b165e4b5f0575e32f92840f28ee7f4c4a73963f) 
		result=$(./sysclient accountAdd  0xA03987F3f20963287ad25e574ab1F9893189Cf3a 0x7e07a3badf7bceeff5b234e5051f088da563e69d9e53051d0a28cc70938536ee) 
		

		if grep -q "succesfully" <<< "$result";
		then
			./system >> logs/ourInitServer.out &
		fi

		;;
	
	(system_host_1) 
		node host_server/host_server.js >> logs/host_server.out &	

		counter=1

		while [ $counter -gt 0 ]
		do
			result=$(nc -z -v localhost 8001 2>&1)
			if grep -q "succeeded" <<< "$result";
			then
				counter=$(($counter - 1))
			fi

		done

		result=$(./sysclient accountAdd  0x7846Cea3FFe81FE4ff84dA127F6Ff25B107D45C5 0x230d3e972fdaef02e991e0d5aaa46f2a28382d41fc76fdae6ce2464c7ad989df)

		echo $result >> logs/sysclient.out

		if grep -q "succesfully" <<< "$result";
		then
			./system >> logs/ourInitServer.out &
		fi
		;;

	(system_host_2)
		node host_server/host_server.js >> logs/host_server.out &

		counter=1

		while [ $counter -gt 0 ]
		do
			result=$(nc -z -v localhost 8001 2>&1)
			if grep -q "succeeded" <<< "$result";
			then
				counter=$(($counter - 1))
			fi
		done

		result=$(./sysclient accountAdd  0x2D9883cc0788bC4Ff6E75C34877DF27E9E22cdD2 0xe47ac19626f743d6ae44f4a0947d0b2b58278ebf8d97da8949326d1a064ed7b4)
		if grep -q "succesfully" <<< "$result";
		then
			./system >> logs/ourInitServer.out &
		fi
		;;

	(system_host_3)
		node host_server/host_server.js >> logs/host_server.out &

		counter=1

		while [ $counter -gt 0 ]
		do
			result=$(nc -z -v localhost 8001 2>&1)
			if grep -q "succeeded" <<< "$result";
			then
				counter=$(($counter - 1))
			fi
		done

		result=$(./sysclient accountAdd  0xF75372C861a5f3Ce2b28Ad0785AeBc0a32a86883 0xc8532c24831cc58d3ace4c213e2a86fd58485215de60527ca64a298a900c8a46)
		if grep -q "succesfully" <<< "$result";
		then
			./system >> logs/ourInitServer.out &
		fi
		;;

	(system_host_4)
		node host_server/host_server.js >> logs/host_server.out &

		counter=1

		while [ $counter -gt 0 ]
		do
			result=$(nc -z -v localhost 8001 2>&1)
			if grep -q "succeeded" <<< "$result";
			then
				counter=$(($counter - 1))
			fi
		done

		result=$(./sysclient accountAdd  0x3c87D59392F0D5df621B950FA090C40a54F372a6 0x44b2854e4e6482963523df4d0be909537e2b36192945bb82295c883ca13f413e)

		if grep -q "succesfully" <<< "$result";
		then
			./system >> logs/ourInitServer.out &
		fi
		;;

	(system_host_5)
		node host_server/host_server.js >> logs/host_server.out &

		counter=1

		while [ $counter -gt 0 ]
		do
			result=$(nc -z -v localhost 8001 2>&1)
			if grep -q "succeeded" <<< "$result";
			then
				counter=$(($counter - 1))
			fi
		done

		result=$(./sysclient accountAdd  0xfdC97A4B2Fc441Fd1A7A5f03eFbdbEf704023291 0xa4d2bc885cf99df88e07b90f18a2d09b3c81d7e12f07e9a97c36cfc1144077eb)

		if grep -q "succesfully" <<< "$result";
		then
			./system >> logs/ourInitServer.out &
		fi
		;;

	(system_host_6)
		node host_server/host_server.js >> logs/host_server.out &

		counter=1

		while [ $counter -gt 0 ]
		do
			result=$(nc -z -v localhost 8001 2>&1)
			if grep -q "succeeded" <<< "$result";
			then
				counter=$(($counter - 1))
			fi
		done


		result=$(./sysclient accountAdd  0x26E63eED9FFb84f7fb8ef1198998C3860c0D13db 0xab6c8ee821c22dfa76953842d342ae8e259077a3ee4c73058c757e8c3a2b46bc)
		if grep -q "succesfully" <<< "$result";
		then
			./system >> logs/ourInitServer.out &
		fi
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

