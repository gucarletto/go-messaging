FROM golang:latest

WORKDIR /go/app

RUN apt-get update && apt-get install librkafka-dev

CMD [ "tail", "-f", "/dev/null" ]