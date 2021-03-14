#!/bin/bash
echo "The script starts now"
#build our image from Dockerfile
docker build -t system_1 .

#create containers from our image
docker create -v /home/fotis/truffle-example/build/contracts:/home/fotis/truffle-example/build/contracts --name system_renter -ti system_1
docker create -v /home/fotis/truffle-example/build/contracts:/home/fotis/truffle-example/build/contracts --name system_host_1 -ti system_1
docker create -v /home/fotis/truffle-example/build/contracts:/home/fotis/truffle-example/build/contracts --name system_host_2 -ti system_1

