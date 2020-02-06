FROM dodasts/centos:7-grid as APP

RUN mkdir /app 
WORKDIR /app

ADD tts-cache /usr/local/bin/tts-cache

ENTRYPOINT ["tts-cache"]