# FROM golang:1.13.7 as BASE

# RUN mkdir /app 
# ADD . /app/ 
# WORKDIR /app 
# #RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -mod vendor -o tts-cache .
# RUN GOOS=linux go build -mod vendor -o tts-cache .

FROM dodasts/tts-cache:base-k8s as CACHED

RUN GOOS=linux go build -mod vendor -o tts-cache .

#FROM alpine as APP
FROM dodasts/centos:7-grid as APP

RUN mkdir /app 
WORKDIR /app

COPY --from=0 /app/tts-cache /usr/local/bin/tts-cache

ENTRYPOINT ["tts-cache"]