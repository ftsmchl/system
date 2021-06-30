#!/bin/bash

#first we must check if all ports are open
for i in {3..8}
do
	counter=1
	while [ $counter -gt 0 ]
	do 
		echo "eimai mesa sto while"
		result=$(nc -z -v 172.17.0.$i 8001 2>&1)
		if grep -q "succeeded" <<< "$result";
		then
			echo "172.17.0.$i is up"
			counter=$(($counter - 1))
		fi
	done
done

curl 172.17.0.2:8089/createContracts &

#i=1

#sleep 29 
echo "pame gia to for loop"

#echo "i is $i"

for i in {3..8}
do
	flag=1
	
	while [ $flag -gt 0 ]
	do
		echo "eimai mesa sto while to deutero"
		sleep 1
		newResult=$(curl -X GET 172.17.0.$i:8001/list -H "Accept : */*")
		if grep -q "koble" <<< "$newResult";
		then
			echo "8a exw sigoura auctions"
			flag=$(( $flag - 1 ))
		fi
	done

	echo "172.17.0.$i"
	curl 172.17.0.$i:8089/findAuctions &
	
done


