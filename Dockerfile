FROM golang:1.12.1

WORKDIR /go/src/app
COPY . /go/src/app

RUN apt-get update && \
	apt-get -y install sqlite3 && \
	rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
	apt-get clean

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["go", "run", "main.go"]