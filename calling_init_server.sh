#!/bin/bash

curl 172.17.0.2:8089/createContracts &

#i=1

sleep 29 
echo "pame gia to for loop"

#echo "i is $i"

for i in {3..8}
do
	#echo "172.17.0.$i"
	curl 172.17.0.$i:8089/findAuctions &
	sleep 3 
	
done


