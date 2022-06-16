FROM golang:1.17

WORKDIR /rancher

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./preinstall.go .

CMD ["go","run","preinstall.go"]