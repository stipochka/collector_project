# test-runner/Dockerfile
FROM golang:1.24.4

WORKDIR /app 

COPY . .

RUN go install github.com/onsi/ginkgo/v2/ginkgo@latest
RUN go mod download

CMD go test -v ./tests/blackbox