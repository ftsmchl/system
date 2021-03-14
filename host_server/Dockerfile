FROM node:14 as builder
WORKDIR node/host_server
COPY package*.json ./
RUN npm install
COPY . .

#RUN apk --no-cache add ca-certificates
