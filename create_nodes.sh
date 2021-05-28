#!/bin/bash
echo "The script starts now"
#build our image from Dockerfile
docker build -t system_1 .

#create containers from our image
docker run --hostname system_renter -p 8089:70 -d -v /home/fotis/truffle-example/build/contracts:/home/fotis/truffle-example/build/contracts --name system_renter -ti system_1 
docker run --hostname system_host_1 -d -v /home/fotis/truffle-example/build/contracts:/home/fotis/truffle-example/build/contracts --name system_host_1 -ti system_1
docker run --hostname system_host_2 -d -v /home/fotis/truffle-example/build/contracts:/home/fotis/truffle-example/build/contracts --name system_host_2 -ti system_1
docker run --hostname system_host_3 -d -v /home/fotis/truffle-example/build/contracts:/home/fotis/truffle-example/build/contracts --name system_host_3 -ti system_1
docker run --hostname system_host_4 -d -v /home/fotis/truffle-example/build/contracts:/home/fotis/truffle-example/build/contracts --name system_host_4 -ti system_1
docker run --hostname system_host_5 -d -v /home/fotis/truffle-example/build/contracts:/home/fotis/truffle-example/build/contracts --name system_host_5 -ti system_1
docker run --hostname system_host_6 -d -v /home/fotis/truffle-example/build/contracts:/home/fotis/truffle-example/build/contracts --name system_host_6 -ti system_1
#docker run --hostname system_host_7 -d -v /home/fotis/truffle-example/build/contracts:/home/fotis/truffle-example/build/contracts --name system_host_7 -ti system_1
#docker run --hostname system_host_8 -d -v /home/fotis/truffle-example/build/contracts:/home/fotis/truffle-example/build/contracts --name system_host_8 -ti system_1

