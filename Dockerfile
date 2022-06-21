FROM golang:1.17

WORKDIR /verrazzano

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY preinstall.go .
COPY platform-operator/manifests/profiles platform-operator/manifests

CMD ["go","run","preinstall.go"]