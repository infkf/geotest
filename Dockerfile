FROM golang:1.17.2

ADD . /go/src/app
# ADD ./go.mod /go/go.mod
WORKDIR /go/src/app
RUN ls
RUN go get .
RUN go install .
ENTRYPOINT geotest

EXPOSE 5000