From golang:1.14.3 as builder
WORKDIR /go/src/system

COPY go.mod .
COPY go.sum .

RUN go mod download 

COPY . .

#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sysd .
RUN CGO_ENABLED=0 GOOS=linux go build -o . ./cmd/sysd
RUN CGO_ENABLED=0 GOOS=linux go build -o . ./cmd/sysclient



#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#WORKDIR /root/
#COPY --from=builder /go/src/system/sysd . 
#COPY --from=builder /go/src/system/sysclient .
#CMD ["gorilla_try"]

FROM ubuntu:latest
WORKDIR /root/
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

RUN apt-get update \
    && apt-get install -y curl \
    && apt-get -y autoclean

RUN apt-get install -y vim 
RUN apt-get install -y net-tools 

ENV NVM_DIR /usr/local/nvm
ENV NODE_VERSION 12.18.3 
RUN curl --silent -o- https://raw.githubusercontent.com/creationix/nvm/v0.31.2/install.sh | bash
#RUN apt-install -y nodejs 
#RUN apt-install -y npm
RUN source $NVM_DIR/nvm.sh \
	&& nvm install $NODE_VERSION \
	&& nvm alias default $NODE_VERSION \
	&& nvm use default
ENV NODE_PATH $NVM_DIR/v$NODE_VERSION/lib/node_modules
ENV PATH $NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH
RUN node -v
RUN npm -v
RUN npm install web3
COPY --from=builder /go/src/system/sysd . 
COPY --from=builder /go/src/system/sysclient .
COPY ./host_server ./host_server
COPY ./renter_server ./renter_server
COPY system ./system
#ENTRYPOINT ./sysd &
#CMD ["./sysd &"]
#RUN ["./sysd", "&"]
