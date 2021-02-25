FROM golang:1.15.8

WORKDIR /go/src/

COPY ./client ./
RUN go mod download

ENTRYPOINT ["go"]
CMD ["test", "-v", "-count=1" ]