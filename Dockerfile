FROM golang:alpine

WORKDIR /go/src/geofinder
ADD . /go/src/geofinder
#Usable For Fetch With Govendor External Packages
#RUN cd /go/src/geofinder && export GOPATH="/go" && apk add --no-cache git mercurial && go get -u github.com/kardianos/govendor && CGO_ENABLED=0 GOOS=linux govendor build -a -installsuffix cgo -o geofinder
RUN cd /go/src/geofinder && export GOPATH="/go" && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o geofinder
ENTRYPOINT ./geofinder

#For Multi Stage Build & Minimize Image Size -> Docker Deamon 17.05
#FROM alpine:latest
#WORKDIR /app
#COPY --from=0 /go/src/geofinder/geofinder /app/
#EXPOSE 8080
#ENTRYPOINT ./geofinder