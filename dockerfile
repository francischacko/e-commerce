FROM golang:1.20.3-alpine3.17
RUN mkdir /app
WORKDIR /app
ADD . /app


COPY go.mod .
COPY go.sum .

ENV GO111MODULE=on
RUN go mod download

RUN go build -o main .
CMD ["/app/main"]


